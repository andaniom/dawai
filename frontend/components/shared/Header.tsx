"use client";

import { useAuthStore } from "@/lib/store";
import { SchoolSwitcher } from "./SchoolSwitcher";

export function Header() {
  const { user } = useAuthStore();

  return (
    <header className="flex items-center justify-between p-4 border-b bg-card">
      <SchoolSwitcher />
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <span>{user?.name}</span>
        <span className="text-xs bg-muted px-2 py-0.5 rounded">
          {user?.roles[0]}
        </span>
      </div>
    </header>
  );
}
