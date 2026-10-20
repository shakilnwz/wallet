import React from 'react';
import { ChevronLeft, ChevronRight, Calendar } from 'lucide-react';
import { formatMonthDisplay } from '../utils/formatters';

interface MonthSelectorProps {
  selectedMonth: string; // 'YYYY-MM' or 'all'
  availableMonths: string[];
  onChange: (month: string) => void;
}

export const MonthSelector: React.FC<MonthSelectorProps> = ({
  selectedMonth,
  availableMonths,
  onChange,
}) => {
  const isAllTime = selectedMonth === 'all';

  const handlePrev = () => {
    if (isAllTime && availableMonths.length > 0) {
      onChange(availableMonths[0]);
      return;
    }
    const idx = availableMonths.indexOf(selectedMonth);
    if (idx !== -1 && idx < availableMonths.length - 1) {
      onChange(availableMonths[idx + 1]);
    } else if (idx === -1 && availableMonths.length > 0) {
      onChange(availableMonths[0]);
    }
  };

  const handleNext = () => {
    if (isAllTime) return;
    const idx = availableMonths.indexOf(selectedMonth);
    if (idx > 0) {
      onChange(availableMonths[idx - 1]);
    }
  };

  const currentIndex = availableMonths.indexOf(selectedMonth);
  const canGoNext = !isAllTime && currentIndex > 0;
  const canGoPrev = !isAllTime && currentIndex < availableMonths.length - 1;

  const displayMonths = Array.from(
    new Set([selectedMonth !== 'all' ? selectedMonth : null, ...availableMonths].filter(Boolean))
  ) as string[];

  return (
    <div className="flex items-center justify-between bg-surface rounded-2xl p-1.5 border border-surface-border">
      <button
        onClick={handlePrev}
        disabled={!canGoPrev && !isAllTime}
        aria-label="Previous Month"
        className="p-2 rounded-xl text-zinc-400 hover:text-white hover:bg-surface-elevated disabled:opacity-30 disabled:pointer-events-none active:scale-95 transition-all"
      >
        <ChevronLeft className="w-5 h-5" />
      </button>

      <div className="flex items-center gap-2">
        <Calendar className="w-4 h-4 text-blue-400" />
        <select
          value={selectedMonth}
          onChange={(e) => onChange(e.target.value)}
          className="bg-transparent text-sm font-semibold text-white focus:outline-none cursor-pointer pr-1"
        >
          {displayMonths.map((m) => (
            <option key={m} value={m} className="bg-surface text-white">
              {formatMonthDisplay(m)}
            </option>
          ))}
          <option value="all" className="bg-surface text-white">
            All Time
          </option>
        </select>
      </div>

      <button
        onClick={handleNext}
        disabled={!canGoNext}
        aria-label="Next Month"
        className="p-2 rounded-xl text-zinc-400 hover:text-white hover:bg-surface-elevated disabled:opacity-30 disabled:pointer-events-none active:scale-95 transition-all"
      >
        <ChevronRight className="w-5 h-5" />
      </button>
    </div>
  );
};
