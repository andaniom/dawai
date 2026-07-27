import Dexie, { Table } from 'dexie';

export interface PendingAssessment {
  id?: number; // Local DB key
  idempotency_key: string; // UUID from frontend
  student_id: string;
  song_id: string;
  scores: Record<string, number>; // {component_id: score}
  feedback?: string;
  created_at: number; // timestamp
  synced: boolean;
  sync_error?: string;
}

export class AppDB extends Dexie {
  pending_assessments!: Table<PendingAssessment>;

  constructor() {
    super('DAWAI');
    this.version(1).stores({
      pending_assessments: '++id, idempotency_key, synced',
    });
  }
}

export const db = new AppDB();
