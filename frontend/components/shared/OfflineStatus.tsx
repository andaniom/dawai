"use client";

import { useOfflineQueue } from "@/hooks";
import { CloudOff, RefreshCw } from "lucide-react";

export function OfflineStatus() {
  const { pending, syncing } = useOfflineQueue();

  if (pending === 0 && !syncing) return null;

  return (
    <div className="fixed bottom-4 right-4 bg-amber-100 border border-amber-300 text-amber-800 px-4 py-2 rounded-md shadow-lg flex items-center space-x-2 z-50">
      {syncing ? (
        <>
          <RefreshCw className="h-4 w-4 animate-spin" />
          <span className="text-sm font-medium">Syncing...</span>
        </>
      ) : (
        <>
          <CloudOff className="h-4 w-4" />
          <span className="text-sm font-medium">
            {pending} Pending sync
          </span>
        </>
      )}
    </div>
  );
}
