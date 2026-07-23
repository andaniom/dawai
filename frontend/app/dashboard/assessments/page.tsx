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
