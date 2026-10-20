import React from 'react';
import { Transaction } from '../types';
import { formatCurrency, formatDateDisplay } from '../utils/formatters';
import { ArrowDownLeft, ArrowUpRight, Building2, Banknote, Trash2, Tag, Edit3 } from 'lucide-react';

interface TransactionListProps {
  transactions: Transaction[];
  onDelete: (id: string) => void;
  onEdit?: (tx: Transaction) => void;
}

export const TransactionList: React.FC<TransactionListProps> = ({
  transactions,
  onDelete,
  onEdit
}) => {
  if (transactions.length === 0) {
    return (
      <div className="bg-surface rounded-3xl p-8 border border-surface-border text-center flex flex-col items-center justify-center min-h-[220px]">
        <div className="w-12 h-12 rounded-2xl bg-zinc-800/80 flex items-center justify-center text-zinc-400 mb-3">
          <Tag className="w-6 h-6" />
        </div>
        <h3 className="text-base font-semibold text-white mb-1">No Transactions Found</h3>
        <p className="text-xs text-zinc-400 max-w-[240px]">
          There are no records for this period. Tap the plus button below to log an expense or income.
        </p>
      </div>
    );
  }

  // Group transactions by date
  const grouped: Record<string, Transaction[]> = {};
  for (const tx of transactions) {
    if (!grouped[tx.date]) grouped[tx.date] = [];
    grouped[tx.date].push(tx);
  }

  const sortedDates = Object.keys(grouped).sort((a, b) => b.localeCompare(a));

  return (
    <div className="space-y-4">
      {sortedDates.map((dateStr) => {
        const txs = grouped[dateStr];
        const dayNet = txs.reduce((sum, tx) => {
          return sum + (tx.type === 'income' ? tx.amount : -tx.amount);
        }, 0);

        return (
          <div key={dateStr} className="space-y-2">
            {/* Date Header */}
            <div className="flex items-center justify-between px-2 text-xs font-semibold text-zinc-400">
              <span>{formatDateDisplay(dateStr)}</span>
              <span className={`text-[11px] font-medium ${dayNet >= 0 ? 'text-emerald-400' : 'text-zinc-400'}`}>
                {dayNet !== 0 ? (dayNet > 0 ? `+${formatCurrency(dayNet)}` : formatCurrency(dayNet)) : ''}
              </span>
            </div>

            {/* Transaction Cards for this Date */}
            <div className="bg-surface rounded-2xl border border-surface-border divide-y divide-surface-border/60 overflow-hidden shadow-sm">
              {txs.map((tx) => (
                <div
                  key={tx.id}
                  className="p-3.5 flex items-center justify-between hover:bg-surface-elevated active:bg-surface-elevated/70 transition-colors"
                >
                  <div className="flex items-center gap-3 min-w-0 flex-1">
                    <div
                      className={`w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0 ${
                        tx.type === 'income'
                          ? 'bg-emerald-500/15 text-emerald-400 border border-emerald-500/25'
                          : 'bg-rose-500/15 text-rose-400 border border-rose-500/25'
                      }`}
                    >
                      {tx.type === 'income' ? (
                        <ArrowDownLeft className="w-4 h-4" />
                      ) : (
                        <ArrowUpRight className="w-4 h-4" />
                      )}
                    </div>

                    <div className="min-w-0 flex-1">
                      <div className="text-sm font-semibold text-white truncate capitalize">
                        {tx.benf}
                      </div>
                      <div className="flex items-center gap-1.5 mt-0.5 text-[11px] text-zinc-400">
                        <span className="capitalize px-1.5 py-0.2 rounded bg-zinc-800 text-zinc-300">
                          {tx.exp_cat}
                        </span>
                        <span>•</span>
                        <span className="flex items-center gap-0.5">
                          {tx.account === 'Account' ? (
                            <Building2 className="w-3 h-3 text-blue-400" />
                          ) : (
                            <Banknote className="w-3 h-3 text-emerald-400" />
                          )}
                          {tx.account}
                        </span>
                        {tx.notes && (
                          <>
                            <span>•</span>
                            <span className="truncate max-w-[120px] text-zinc-500">{tx.notes}</span>
                          </>
                        )}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 pl-3">
                    <div className="text-right">
                      <div
                        className={`text-sm font-bold tracking-tight ${
                          tx.type === 'income' ? 'text-emerald-400' : 'text-zinc-100'
                        }`}
                      >
                        {tx.type === 'income' ? '+' : '-'}
                        {formatCurrency(tx.amount)}
                      </div>
                    </div>

                    <div className="flex items-center">
                      {onEdit && (
                        <button
                          onClick={() => onEdit(tx)}
                          title="Edit transaction"
                          className="p-1.5 rounded-lg text-zinc-500 hover:text-blue-400 active:scale-95 transition-all"
                        >
                          <Edit3 className="w-3.5 h-3.5" />
                        </button>
                      )}
                      <button
                        onClick={() => {
                          if (confirm(`Delete "${tx.benf}"?`)) {
                            onDelete(tx.id);
                          }
                        }}
                        title="Delete transaction"
                        className="p-1.5 rounded-lg text-zinc-500 hover:text-rose-400 active:scale-95 transition-all"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
};
