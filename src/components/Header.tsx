import React from 'react';
import { Wallet, Database, RefreshCw } from 'lucide-react';

interface HeaderProps {
  onOpenSync: () => void;
  onRefresh: () => void;
  isSaving?: boolean;
}

export const Header: React.FC<HeaderProps> = ({ onOpenSync, onRefresh, isSaving }) => {
  return (
    <header className="sticky top-0 z-30 bg-surface/90 backdrop-blur-md border-b border-surface-border px-4 py-3 flex items-center justify-between pt-safe">
      <div className="flex items-center gap-2.5">
        <div className="w-9 h-9 rounded-xl bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400">
          <Wallet className="w-5 h-5" />
        </div>
        <div>
          <div className="flex items-center gap-1.5">
            <h1 className="text-base font-bold text-white tracking-tight">Wallet</h1>
            <span className="text-[10px] uppercase tracking-wider font-semibold px-1.5 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              SQLite
            </span>
          </div>
          <p className="text-[11px] text-zinc-400">Personal Finance</p>
        </div>
      </div>

      <div className="flex items-center gap-1">
        {isSaving ? (
          <span className="text-[11px] text-zinc-400 flex items-center gap-1 px-2 py-1">
            <RefreshCw className="w-3 h-3 animate-spin text-blue-400" />
            Saving...
          </span>
        ) : null}
        <button
          onClick={onRefresh}
          title="Refresh view"
          className="p-2 rounded-lg text-zinc-400 hover:text-white hover:bg-surface-elevated active:scale-95 transition-all"
        >
          <RefreshCw className="w-4 h-4" />
        </button>
        <button
          onClick={onOpenSync}
          title="Database & Sync"
          className="p-2 rounded-lg text-zinc-400 hover:text-white hover:bg-surface-elevated active:scale-95 transition-all flex items-center gap-1"
        >
          <Database className="w-4 h-4 text-zinc-300" />
        </button>
      </div>
    </header>
  );
};
