"use client";

import { useAuthStore } from "@/lib/store";

export function SchoolSwitcher() {
  const { schoolId } = useAuthStore();
  // ponytail: single school for MVP, dropdown when multi-school
  return (
    <div className="flex items-center gap-2">
      <span className="font-medium text-sm">Current School</span>
      <span className="text-sm text-muted-foreground">
        {schoolId || "Loading..."}
      </span>
    </div>
  );
}
