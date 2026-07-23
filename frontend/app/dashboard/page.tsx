"use client";

import { useAuthStore } from "@/lib/store";

export default function DashboardPage() {
  const { user } = useAuthStore();

  return (
    <div>
      <h1 className="text-2xl font-bold">Welcome, {user?.name}</h1>
      <p className="text-muted-foreground mt-2">
        Select a subject and student to begin assessing.
      </p>
    </div>
  );
}
