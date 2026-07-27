"use client";

import Link from "next/link";
import { useAuthStore } from "@/lib/store";

export function Sidebar() {
  const { user, clearAuth } = useAuthStore();

  return (
    <aside className="w-64 border-r bg-card p-4 flex flex-col">
      <h2 className="text-lg font-semibold mb-4">DAWAI</h2>
      <nav className="flex-1 space-y-2">
        <Link href="/dashboard" className="block p-2 hover:bg-accent rounded">
          Dashboard
        </Link>
        <Link
          href="/dashboard/assessments"
          className="block p-2 hover:bg-accent rounded"
        >
          Assessments
        </Link>
        <Link
          href="/dashboard/results"
          className="block p-2 hover:bg-accent rounded"
        >
          Results
        </Link>
        {user?.roles.includes("school_admin") && (
          <>
            <Link
              href="/dashboard/subjects"
              className="block p-2 hover:bg-accent rounded"
            >
              Subjects
            </Link>
            <Link
              href="/dashboard/users"
              className="block p-2 hover:bg-accent rounded"
            >
              Users
            </Link>
          </>
        )}
      </nav>
      <button
        onClick={clearAuth}
        className="text-sm text-destructive hover:underline mt-4"
      >
        Sign Out
      </button>
    </aside>
  );
}
