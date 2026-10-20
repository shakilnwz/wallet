import React from 'react';
import { Home, Receipt, BarChart3, Database, Plus } from 'lucide-react';

export type NavTab = 'dashboard' | 'ledger' | 'analytics' | 'sync';

interface BottomNavProps {
  currentTab: NavTab;
  onTabChange: (tab: NavTab) => void;
  onOpenAdd: () => void;
}

export const BottomNav: React.FC<BottomNavProps> = ({
  currentTab,
  onTabChange,
  onOpenAdd
}) => {
  return (
    <nav className="fixed bottom-0 left-0 right-0 z-40 bg-surface/95 backdrop-blur-lg border-t border-surface-border px-3 pb-safe pt-2">
      <div className="max-w-md mx-auto flex items-center justify-between relative">
        {/* Left Side: Dashboard & Ledger */}
        <div className="flex items-center gap-1 flex-1 justify-around pr-8">
          <button
            onClick={() => onTabChange('dashboard')}
            className={`flex flex-col items-center gap-1 py-1 px-3 rounded-xl transition-all active:scale-95 ${
              currentTab === 'dashboard'
                ? 'text-blue-400 font-bold'
                : 'text-zinc-400 hover:text-zinc-200 font-medium'
            }`}
          >
            <Home className="w-5 h-5" />
            <span className="text-[10px]">Home</span>
          </button>

          <button
            onClick={() => onTabChange('ledger')}
            className={`flex flex-col items-center gap-1 py-1 px-3 rounded-xl transition-all active:scale-95 ${
              currentTab === 'ledger'
                ? 'text-blue-400 font-bold'
                : 'text-zinc-400 hover:text-zinc-200 font-medium'
            }`}
          >
            <Receipt className="w-5 h-5" />
            <span className="text-[10px]">Ledger</span>
          </button>
        </div>

        {/* Center Prominent Add (+) Button */}
        <div className="absolute left-1/2 -top-5 -translate-x-1/2">
          <button
            onClick={onOpenAdd}
            aria-label="Add transaction"
            className="w-13 h-13 w-[52px] h-[52px] rounded-full bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-600/30 flex items-center justify-center active:scale-90 transition-all border-4 border-background"
          >
            <Plus className="w-6 h-6 stroke-[2.5]" />
          </button>
        </div>

        {/* Right Side: Analytics & Sync */}
        <div className="flex items-center gap-1 flex-1 justify-around pl-8">
          <button
            onClick={() => onTabChange('analytics')}
            className={`flex flex-col items-center gap-1 py-1 px-3 rounded-xl transition-all active:scale-95 ${
              currentTab === 'analytics'
                ? 'text-blue-400 font-bold'
                : 'text-zinc-400 hover:text-zinc-200 font-medium'
            }`}
          >
            <BarChart3 className="w-5 h-5" />
            <span className="text-[10px]">Insights</span>
          </button>

          <button
            onClick={() => onTabChange('sync')}
            className={`flex flex-col items-center gap-1 py-1 px-3 rounded-xl transition-all active:scale-95 ${
              currentTab === 'sync'
                ? 'text-blue-400 font-bold'
                : 'text-zinc-400 hover:text-zinc-200 font-medium'
            }`}
          >
            <Database className="w-5 h-5" />
            <span className="text-[10px]">Sync & DB</span>
          </button>
        </div>
      </div>
    </nav>
  );
};
