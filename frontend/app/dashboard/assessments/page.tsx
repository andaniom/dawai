"use client";

import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import apiClient from "@/lib/api";
import { ApiResponse, Student, Subject } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { AssessmentModal } from "@/components/assessment/AssessmentModal";
import { useLiveQuery } from "dexie-react-hooks";
import { db } from "@/lib/db";
import { Clock } from "lucide-react";

export default function AssessmentsPage() {
  const [selectedSubject, setSelectedSubject] = useState<string | null>(null);
  const [selectedStudent, setSelectedStudent] = useState<{ id: string; name: string } | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const { data: subjectsData } = useQuery({
    queryKey: ["subjects"],
    queryFn: () => apiClient.get<ApiResponse<Subject[]>>("/api/subjects"),
  });

  const { data: studentsData, isLoading } = useQuery({
    queryKey: ["students"],
    queryFn: () => apiClient.get<ApiResponse<Student[]>>("/api/students"),
  });

  const subjects: Subject[] = subjectsData?.data?.data ?? [];
  const students: Student[] = studentsData?.data?.data ?? [];

  const pendingAssessments = useLiveQuery(
    () => db.pending_assessments.where("synced").equals(0).toArray()
  );

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Assessments</h1>

      <div className="w-64">
        <Select value={selectedSubject ?? ""} onValueChange={setSelectedSubject}>
          <SelectTrigger>
            <SelectValue placeholder="Select subject" />
          </SelectTrigger>
          <SelectContent>
            {subjects.map((s: Subject) => (
              <SelectItem key={s.id} value={s.id}>
                {s.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {isLoading ? (
        <p>Loading students...</p>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Student</TableHead>
              <TableHead>Class</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {students.map((s: Student) => (
              <TableRow key={s.id}>
                <TableCell className="font-medium">{s.name}</TableCell>
                <TableCell>{s.class}</TableCell>
                <TableCell>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={!selectedSubject}
                    onClick={() => {
                      setSelectedStudent({ id: s.id, name: s.name });
                      setModalOpen(true);
                    }}
                  >
                    Assess
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {students.length === 0 && (
              <TableRow>
                <TableCell colSpan={3} className="text-center text-muted-foreground py-8">
                  No students found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      )}

      {pendingAssessments && pendingAssessments.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <Clock className="w-5 h-5 text-amber-500" />
            Pending Assessments
          </h2>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Student</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pendingAssessments.map((p) => {
                const studentName = students.find((s) => s.id === p.student_id)?.name ?? "Unknown";
                const subjectName = subjects.find((s) => s.id === p.song_id)?.name ?? "Unknown";
                return (
                  <TableRow key={p.idempotency_key}>
                    <TableCell>{studentName}</TableCell>
                    <TableCell>{subjectName}</TableCell>
                    <TableCell>
                      <span className="inline-flex items-center rounded-full bg-amber-50 px-2 py-1 text-xs font-medium text-amber-700 ring-1 ring-inset ring-amber-600/20">
                        Pending sync
                      </span>
                      {p.sync_error && (
                        <span className="ml-2 text-xs text-destructive">
                          {p.sync_error}
                        </span>
                      )}
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      )}

      {selectedStudent && selectedSubject && (
        <AssessmentModal
          open={modalOpen}
          onOpenChange={setModalOpen}
          studentId={selectedStudent.id}
          studentName={selectedStudent.name}
          subjectId={selectedSubject}
          subjectName={subjects.find((s: Subject) => s.id === selectedSubject)?.name ?? ""}
        />
      )}
    </div>
  );
}
