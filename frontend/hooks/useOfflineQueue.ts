'use client';

import { useEffect, useState, useCallback } from 'react';
import { db } from '@/lib/db';
import { v4 as uuidv4 } from 'uuid';
import apiClient from '@/lib/api';

export interface QueuedAssessmentResult {
  queued: boolean;
  idempotency_key: string;
  id?: number;
}

export function useOfflineQueue() {
  const [pending, setPending] = useState(0);
  const [syncing, setSyncing] = useState(false);

  // On mount, count pending assessments
  useEffect(() => {
    const updateCount = async () => {
      const count = await db.pending_assessments.where('synced').equals(0).count();
      setPending(count);
    };
    updateCount();
  }, []);

  // Listen for network online
  useEffect(() => {
    const handleOnline = async () => {
      setSyncing(true);
      await syncPendingAssessments();
      setSyncing(false);
    };

    window.addEventListener('online', handleOnline);
    return () => window.removeEventListener('online', handleOnline);
  }, []);

  // Queue assessment submission
  const queueAssessment = useCallback(
    async (data: {
      student_id: string;
      subject_id: string;
      scores: { rubric_component_id: string; score: number }[];
      feedback?: string;
    }): Promise<QueuedAssessmentResult> => {
      const idempotency_key = uuidv4();

      try {
        // Try online submission
        const response = await apiClient.post('/api/assessments', data, {
          headers: { 'Idempotency-Key': idempotency_key },
        });
        if (response.status === 200 || response.status === 201) {
          return {
            queued: false,
            idempotency_key,
          };
        }
        throw new Error('Network error');
      } catch (error: any) {
        // Save to queue if offline or error
        const id = await db.pending_assessments.add({
          idempotency_key,
          student_id: data.student_id,
          song_id: data.subject_id, // subject_id maps to song_id in schema
          scores: data.scores.reduce(
            (acc, item) => {
              acc[item.rubric_component_id] = item.score;
              return acc;
            },
            {} as Record<string, number>
          ),
          feedback: data.feedback,
          created_at: Date.now(),
          synced: false,
        });

        setPending((p) => p + 1);
        return {
          queued: true,
          idempotency_key,
          id,
        };
      }
    },
    []
  );

  // Sync when online
  const syncPendingAssessments = useCallback(async () => {
    const pending_list = await db.pending_assessments
      .where('synced')
      .equals(0)
      .toArray();

    let synced_count = 0;

    for (const item of pending_list) {
      try {
        // Reconstruct data format
        const scores = Object.entries(item.scores).map(([rubric_component_id, score]) => ({
          rubric_component_id,
          score,
        }));

        const response = await apiClient.post(
          '/api/assessments',
          {
            student_id: item.student_id,
            subject_id: item.song_id,
            scores,
            feedback: item.feedback,
          },
          {
            headers: { 'Idempotency-Key': item.idempotency_key },
          }
        );

        if (response.status === 200 || response.status === 201) {
          if (item.id !== undefined) {
            await db.pending_assessments.update(item.id, { synced: true });
            synced_count++;
          }
        }
      } catch (error: any) {
        if (item.id !== undefined) {
          await db.pending_assessments.update(item.id, {
            sync_error: error.message,
          });
        }
      }
    }

    // Update pending count
    const remaining = await db.pending_assessments.where('synced').equals(0).count();
    setPending(remaining);
  }, []);

  return { queueAssessment, pending, syncing, syncPendingAssessments };
}
