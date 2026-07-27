"use client";

import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { Button } from "@/components/ui/button";
import { utils, writeFile } from "xlsx";
import { Assessment, ApiResponse } from "@/lib/types";

interface AssessmentWithComponents extends Assessment {
  subject_name?: string;
  student_name?: string;
  components?: { name: string; score: number }[];
}

export default function ReportsPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["all-assessments-export"],
    queryFn: () => apiClient.get<ApiResponse<AssessmentWithComponents[]>>("/api/assessments"),
  });

  const exportKurMer = () => {
    if (!data?.data?.data) return;
    
    const exportData = data.data.data.map(a => {
      const base: any = {
        "Date": new Date(a.submitted_at).toLocaleDateString(),
        "Subject": a.subject_name || "Unknown",
        "Student": a.student_name || "Unknown",
        "Feedback": a.feedback || "",
      };
      
      a.components?.forEach(c => {
        base[c.name] = c.score;
      });
      
      return base;
    });

    const ws = utils.json_to_sheet(exportData);
    const wb = utils.book_new();
    utils.book_append_sheet(wb, ws, "KurMer");
    writeFile(wb, "KurMer_Report.xlsx");
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Reports</h1>
      <div className="p-6 border rounded-lg bg-card">
        <h2 className="text-xl mb-4">Export Data</h2>
        <p className="text-muted-foreground mb-4">Download KurMer Excel report for all assessments.</p>
        <Button onClick={exportKurMer} disabled={isLoading || !data?.data?.data?.length}>
          {isLoading ? "Loading..." : "Download KurMer"}
        </Button>
      </div>
    </div>
  );
}
