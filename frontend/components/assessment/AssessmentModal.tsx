"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { ApiResponse, RubricComponent } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";

interface AssessmentModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  studentId: string;
  studentName: string;
  subjectId: string;
  subjectName: string;
}

export function AssessmentModal({
  open,
  onOpenChange,
  studentId,
  studentName,
  subjectId,
  subjectName,
}: AssessmentModalProps) {
  const queryClient = useQueryClient();

  const { data: rubricData } = useQuery({
    queryKey: ["rubric", subjectId],
    queryFn: () =>
      apiClient.get<ApiResponse<RubricComponent[]>>(`/api/subjects/${subjectId}/rubric`),
    enabled: open,
  });

  const rubric: RubricComponent[] = rubricData?.data?.data ?? [];

  const createMutation = useMutation({
    mutationFn: (data: {
      student_id: string;
      subject_id: string;
      scores: { rubric_component_id: string; score: number }[];
      feedback: string;
    }) => apiClient.post("/api/assessments", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["assessments"] });
      onOpenChange(false);
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const scores = rubric.map((comp) => ({
      rubric_component_id: comp.id,
      score: Number(formData.get(`score_${comp.id}`)) || 0,
    }));
    const feedback = formData.get("feedback") as string;

    createMutation.mutate({
      student_id: studentId,
      subject_id: subjectId,
      scores,
      feedback,
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>
            Assess {studentName} — {subjectName}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          {rubric.map((comp) => (
            <div key={comp.id} className="space-y-2">
              <Label htmlFor={`score_${comp.id}`}>
                {comp.name} (scale {comp.scale_min}–{comp.scale_max}, weight {comp.weight})
              </Label>
              <Input
                id={`score_${comp.id}`}
                name={`score_${comp.id}`}
                type="number"
                min={comp.scale_min}
                max={comp.scale_max}
                required
                defaultValue={comp.scale_min}
              />
            </div>
          ))}
          <div className="space-y-2">
            <Label htmlFor="feedback">Feedback (optional)</Label>
            <Textarea id="feedback" name="feedback" placeholder="Notes for the student..." />
          </div>
          <Button type="submit" className="w-full" disabled={createMutation.isPending}>
            {createMutation.isPending ? "Submitting..." : "Submit Assessment"}
          </Button>
          {createMutation.isError && (
            <p className="text-sm text-destructive">
              {(createMutation.error as any)?.response?.data?.error?.message || "Failed to submit assessment"}
            </p>
          )}
        </form>
      </DialogContent>
    </Dialog>
  );
}
