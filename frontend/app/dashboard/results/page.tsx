"use client";

import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { ApiResponse, Assessment } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

interface AssessmentWithComponents extends Assessment {
  subject_name?: string;
  components?: { name: string; score: number }[];
}

export default function ResultsPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["my-assessments"],
    queryFn: () =>
      apiClient.get<ApiResponse<AssessmentWithComponents[]>>("/api/assessments"),
  });

  const assessments: AssessmentWithComponents[] = data?.data?.data ?? [];

  if (isLoading) return <p className="p-6">Loading results...</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">My Results</h1>

      {assessments.length === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">
            No assessments yet.
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {assessments.map((a: AssessmentWithComponents) => (
            <Card key={a.id}>
              <CardHeader>
                <CardTitle className="text-lg">
                  {a.subject_name ?? "Subject"} —{" "}
                  {new Date(a.submitted_at).toLocaleDateString()}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Component</TableHead>
                      <TableHead>Score</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {a.components?.map((c) => (
                      <TableRow key={c.name}>
                        <TableCell>{c.name}</TableCell>
                        <TableCell>{c.score}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
                {a.feedback && (
                  <p className="text-sm text-muted-foreground mt-4">
                    <strong>Feedback:</strong> {a.feedback}
                  </p>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
