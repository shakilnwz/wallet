import React from 'react';
import { ArrowUpRight, ArrowDownLeft, Building2, Banknote, CheckCircle2, AlertTriangle } from 'lucide-react';
import { AccountSummary } from '../types';
import { formatCurrency } from '../utils/formatters';

interface BalanceCardProps {
  summary: AccountSummary;
  monthLabel: string;
}

export const BalanceCard: React.FC<BalanceCardProps> = ({ summary, monthLabel }) => {
  return (
    <div className="bg-gradient-to-br from-surface to-surface-elevated rounded-3xl p-5 border border-surface-border shadow-xl">
      <div className="flex items-center justify-between text-zinc-400 text-xs font-medium uppercase tracking-wider mb-1">
        <span>Net Balance ({monthLabel})</span>
        {summary.isBalanced ? (
          <span className="flex items-center gap-1 text-emerald-400 font-semibold lowercase text-[11px] bg-emerald-500/10 px-2 py-0.5 rounded-full border border-emerald-500/20">
            <CheckCircle2 className="w-3 h-3" />
            balanced
          </span>
        ) : (
          <span className="flex items-center gap-1 text-amber-400 font-semibold lowercase text-[11px] bg-amber-500/10 px-2 py-0.5 rounded-full border border-amber-500/20">
            <AlertTriangle className="w-3 h-3" />
            reconcile
          </span>
        )}
      </div>

      <div className="text-3xl font-extrabold tracking-tight text-white my-2">
        {formatCurrency(summary.totalBalance)}
      </div>

      {/* Income and Expense Grid */}
      <div className="grid grid-cols-2 gap-3 mt-4 pt-4 border-t border-surface-border/60">
        <div className="bg-background/40 rounded-2xl p-3 border border-surface-border/40">
          <div className="flex items-center gap-1.5 text-emerald-400 text-xs font-semibold mb-1">
            <div className="w-5 h-5 rounded-full bg-emerald-500/20 flex items-center justify-center">
              <ArrowDownLeft className="w-3 h-3" />
            </div>
            Income
          </div>
          <div className="text-base font-bold text-emerald-400 truncate">
            {formatCurrency(summary.totalIncome)}
          </div>
        </div>

        <div className="bg-background/40 rounded-2xl p-3 border border-surface-border/40">
          <div className="flex items-center gap-1.5 text-rose-400 text-xs font-semibold mb-1">
            <div className="w-5 h-5 rounded-full bg-rose-500/20 flex items-center justify-center">
              <ArrowUpRight className="w-3 h-3" />
            </div>
            Expense
          </div>
          <div className="text-base font-bold text-rose-400 truncate">
            {formatCurrency(summary.totalExpense)}
          </div>
        </div>
      </div>

      {/* Accounts Split: Bank vs Cash & Need Balancing */}
      <div className="mt-4 pt-3 border-t border-surface-border/60 flex items-center justify-between text-xs">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1.5">
            <Building2 className="w-3.5 h-3.5 text-blue-400" />
            <span className="text-zinc-400">Account:</span>
            <span className="font-semibold text-zinc-200">{formatCurrency(summary.bankBalance)}</span>
          </div>

          <div className="flex items-center gap-1.5">
            <Banknote className="w-3.5 h-3.5 text-emerald-400" />
            <span className="text-zinc-400">Cash:</span>
            <span className="font-semibold text-zinc-200">{formatCurrency(summary.cashBalance)}</span>
          </div>
        </div>

        {!summary.isBalanced && (
          <div className="text-[11px] text-amber-400 font-medium">
            Diff: {formatCurrency(summary.needBalancing)}
          </div>
        )}
      </div>
    </div>
  );
};
