import React from 'react';
import { MonthlySummary, AccountSummary } from '../types';
import { formatCurrency, formatMonthDisplay } from '../utils/formatters';
import { PieChart, Building2, Banknote } from 'lucide-react';

interface AnalyticsViewProps {
  monthlySummary: MonthlySummary;
  overallSummary: AccountSummary;
  month: string;
}

export const AnalyticsView: React.FC<AnalyticsViewProps> = ({
  monthlySummary,
  overallSummary,
  month
}) => {
  const isAllTime = month === 'all';
  const income = isAllTime ? overallSummary.totalIncome : monthlySummary.income;
  const expense = isAllTime ? overallSummary.totalExpense : monthlySummary.expense;

  const savingsRate =
    income > 0 ? Math.max(0, Math.round(((income - expense) / income) * 100)) : 0;

  return (
    <div className="space-y-4">
      {/* Monthly/Period Highlights Card */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <PieChart className="w-4 h-4 text-blue-400" />
            <h3 className="text-sm font-bold text-white">
              Spending Analytics ({isAllTime ? 'All Time' : formatMonthDisplay(month)})
            </h3>
          </div>
          <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
            {savingsRate}% saved
          </span>
        </div>

        {/* Breakdown Progress Bars */}
        <div className="space-y-3">
          {monthlySummary.categoryBreakdown.length === 0 ? (
            <p className="text-xs text-zinc-500 text-center py-6">
              No expense categories recorded for this period.
            </p>
          ) : (
            monthlySummary.categoryBreakdown.map((item) => (
              <div key={item.category} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium text-zinc-300 capitalize">{item.category}</span>
                  <div className="flex items-center gap-2">
                    <span className="font-semibold text-white">{formatCurrency(item.amount)}</span>
                    <span className="text-[11px] text-zinc-500 w-9 text-right">{item.percentage}%</span>
                  </div>
                </div>
                <div className="w-full h-2 bg-zinc-800 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-rose-500 to-amber-500 rounded-full transition-all duration-500"
                    style={{ width: `${Math.min(100, Math.max(item.percentage, 3))}%` }}
                  />
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Account Distribution Card */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border">
        <h3 className="text-sm font-bold text-white mb-3">Balance Holdings</h3>
        <div className="grid grid-cols-2 gap-3">
          <div className="p-3.5 rounded-2xl bg-surface-elevated border border-surface-border">
            <div className="flex items-center gap-1.5 text-xs text-blue-400 font-semibold mb-1">
              <Building2 className="w-3.5 h-3.5" />
              Account (Bank/MFS)
            </div>
            <div className="text-base font-bold text-white">
              {formatCurrency(overallSummary.bankBalance)}
            </div>
            <div className="text-[10px] text-zinc-500 mt-1">
              {overallSummary.totalBalance > 0
                ? `${Math.round((overallSummary.bankBalance / overallSummary.totalBalance) * 100)}% of total`
                : ''}
            </div>
          </div>

          <div className="p-3.5 rounded-2xl bg-surface-elevated border border-surface-border">
            <div className="flex items-center gap-1.5 text-xs text-emerald-400 font-semibold mb-1">
              <Banknote className="w-3.5 h-3.5" />
              Cash on Hand
            </div>
            <div className="text-base font-bold text-white">
              {formatCurrency(overallSummary.cashBalance)}
            </div>
            <div className="text-[10px] text-zinc-500 mt-1">
              {overallSummary.totalBalance > 0
                ? `${Math.round((overallSummary.cashBalance / overallSummary.totalBalance) * 100)}% of total`
                : ''}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
