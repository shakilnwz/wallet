import React, { useState, useEffect } from 'react';
import { X, Building2, Banknote, Check, Plus } from 'lucide-react';
import { Transaction, TransactionType, AccountType } from '../types';
import { getTodayDateString } from '../utils/formatters';

interface AddTransactionSheetProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: {
    id?: string;
    date: string;
    type: TransactionType;
    benf: string;
    exp_cat: string;
    amount: number;
    account: AccountType;
    notes?: string;
  }) => void;
  recentCategories: string[];
  recentBeneficiaries: string[];
  editingTransaction?: Transaction | null;
}

const DEFAULT_EXPENSE_CATEGORIES = [
  'Food & Dining',
  'Housing & Rent',
  'Internet',
  'Utilities',
  'Groceries',
  'Transport',
  'Healthcare',
  'Shopping',
  'Personal',
  'Other'
];

const DEFAULT_INCOME_CATEGORIES = [
  'Salary',
  'Carry Forward',
  'Freelance',
  'Investments',
  'Bonus',
  'Gift',
  'Other'
];

export const AddTransactionSheet: React.FC<AddTransactionSheetProps> = ({
  isOpen,
  onClose,
  onSubmit,
  recentCategories,
  recentBeneficiaries,
  editingTransaction
}) => {
  const [type, setType] = useState<TransactionType>('expense');
  const [amountStr, setAmountStr] = useState<string>('');
  const [benf, setBenf] = useState<string>('');
  const [category, setCategory] = useState<string>('');
  const [account, setAccount] = useState<AccountType>('Account');
  const [date, setDate] = useState<string>(getTodayDateString());
  const [notes, setNotes] = useState<string>('');
  const [customCatInput, setCustomCatInput] = useState<boolean>(false);

  useEffect(() => {
    if (editingTransaction) {
      setType(editingTransaction.type);
      setAmountStr(String(editingTransaction.amount));
      setBenf(editingTransaction.benf);
      setCategory(editingTransaction.exp_cat);
      setAccount(editingTransaction.account);
      setDate(editingTransaction.date);
      setNotes(editingTransaction.notes || '');
    } else {
      setType('expense');
      setAmountStr('');
      setBenf('');
      setCategory('');
      setAccount('Account');
      setDate(getTodayDateString());
      setNotes('');
      setCustomCatInput(false);
    }
  }, [editingTransaction, isOpen]);

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const parsedAmount = parseFloat(amountStr);
    if (isNaN(parsedAmount) || parsedAmount <= 0) {
      alert('Please enter a valid amount greater than 0');
      return;
    }

    const finalBenf = benf.trim() || category.trim() || (type === 'income' ? 'Income' : 'Expense');
    const finalCategory = category.trim() || 'other';

    onSubmit({
      id: editingTransaction ? editingTransaction.id : undefined,
      date,
      type,
      benf: finalBenf,
      exp_cat: finalCategory,
      amount: Math.round(parsedAmount * 100) / 100,
      account,
      notes: notes.trim() || undefined
    });

    onClose();
  };

  const handleQuickAddAmount = (addValue: number) => {
    const current = parseFloat(amountStr) || 0;
    setAmountStr(String(current + addValue));
  };

  const categoriesList = Array.from(
    new Set([
      ...(type === 'income' ? DEFAULT_INCOME_CATEGORIES : DEFAULT_EXPENSE_CATEGORIES),
      ...recentCategories
    ])
  );

  const setYesterday = () => {
    const d = new Date();
    d.setDate(d.getDate() - 1);
    const yr = d.getFullYear();
    const mo = String(d.getMonth() + 1).padStart(2, '0');
    const da = String(d.getDate()).padStart(2, '0');
    setDate(`${yr}-${mo}-${da}`);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/75 backdrop-blur-sm animate-in fade-in duration-200">
      <div
        className="w-full sm:max-w-lg bg-surface border-t sm:border border-surface-border sm:rounded-3xl rounded-t-3xl max-h-[92dvh] overflow-y-auto pb-safe shadow-2xl flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Handle bar for mobile */}
        <div className="w-12 h-1.5 bg-zinc-700/60 rounded-full mx-auto mt-3 mb-1 sm:hidden" />

        {/* Header */}
        <div className="px-5 py-3 border-b border-surface-border flex items-center justify-between">
          <h2 className="text-base font-bold text-white">
            {editingTransaction ? 'Edit Transaction' : 'Quick Add Transaction'}
          </h2>
          <button
            onClick={onClose}
            className="p-1.5 rounded-full text-zinc-400 hover:text-white hover:bg-surface-elevated active:scale-90 transition-all"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-5 space-y-4">
          {/* Segmented Control: Expense vs Income */}
          <div className="grid grid-cols-2 p-1 bg-surface-elevated rounded-2xl border border-surface-border">
            <button
              type="button"
              onClick={() => {
                setType('expense');
                if (type === 'income') setCategory('');
              }}
              className={`py-2.5 rounded-xl text-sm font-bold transition-all ${
                type === 'expense'
                  ? 'bg-rose-600 text-white shadow-md'
                  : 'text-zinc-400 hover:text-white'
              }`}
            >
              Expense
            </button>
            <button
              type="button"
              onClick={() => {
                setType('income');
                if (type === 'expense') setCategory('');
              }}
              className={`py-2.5 rounded-xl text-sm font-bold transition-all ${
                type === 'income'
                  ? 'bg-emerald-600 text-white shadow-md'
                  : 'text-zinc-400 hover:text-white'
              }`}
            >
              Income
            </button>
          </div>

          {/* Amount Display */}
          <div className="bg-surface-elevated rounded-2xl p-3.5 border border-surface-border focus-within:border-blue-500 transition-colors">
            <label className="text-[11px] font-semibold uppercase tracking-wider text-zinc-400 block mb-1">
              Amount (BDT)
            </label>
            <div className="flex items-center gap-2">
              <span className="text-2xl font-black text-zinc-300">৳</span>
              <input
                type="number"
                step="any"
                inputMode="decimal"
                value={amountStr}
                onChange={(e) => setAmountStr(e.target.value)}
                placeholder="0.00"
                autoFocus
                required
                className="w-full bg-transparent text-3xl font-extrabold text-white placeholder-zinc-600 focus:outline-none"
              />
            </div>

            {/* Quick Increment Chips */}
            <div className="flex items-center gap-2 mt-3 pt-2.5 border-t border-surface-border/60 overflow-x-auto pb-1 text-xs">
              {[100, 500, 1000, 5000].map((amt) => (
                <button
                  key={amt}
                  type="button"
                  onClick={() => handleQuickAddAmount(amt)}
                  className="px-2.5 py-1 rounded-lg bg-surface border border-surface-border text-zinc-300 hover:text-white hover:bg-surface-border active:scale-95 transition-all flex items-center gap-0.5 flex-shrink-0 font-medium"
                >
                  <Plus className="w-3 h-3 text-zinc-400" />
                  ৳{amt}
                </button>
              ))}
            </div>
          </div>

          {/* Category Chips Selection */}
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <label className="text-xs font-semibold text-zinc-400">Category</label>
              <button
                type="button"
                onClick={() => setCustomCatInput(!customCatInput)}
                className="text-[11px] text-blue-400 hover:underline"
              >
                {customCatInput ? 'Select from list' : '+ Custom category'}
              </button>
            </div>

            {customCatInput ? (
              <input
                type="text"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                placeholder="Enter custom category..."
                className="w-full bg-surface-elevated rounded-xl p-3 text-sm text-white border border-surface-border focus:outline-none focus:border-blue-500"
              />
            ) : (
              <div className="flex flex-wrap gap-1.5 max-h-28 overflow-y-auto pr-1">
                {categoriesList.map((cat) => {
                  const isSelected = category.toLowerCase() === cat.toLowerCase();
                  return (
                    <button
                      key={cat}
                      type="button"
                      onClick={() => {
                        setCategory(cat);
                        if (!benf) setBenf(cat);
                      }}
                      className={`px-3 py-1.5 rounded-xl text-xs font-medium capitalize transition-all active:scale-95 flex items-center gap-1 ${
                        isSelected
                          ? type === 'income'
                            ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/50'
                            : 'bg-rose-500/20 text-rose-300 border border-rose-500/50'
                          : 'bg-surface-elevated text-zinc-300 border border-surface-border hover:bg-zinc-800'
                      }`}
                    >
                      {isSelected && <Check className="w-3 h-3" />}
                      {cat}
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          {/* Beneficiary / Description */}
          <div>
            <label className="text-xs font-semibold text-zinc-400 block mb-1">
              Beneficiary / Payee / Description
            </label>
            <input
              type="text"
              value={benf}
              onChange={(e) => setBenf(e.target.value)}
              placeholder="e.g. Lunch, Wifi bill, Rent, Groceries..."
              list="recent-beneficiaries"
              className="w-full bg-surface-elevated rounded-xl p-3 text-sm text-white border border-surface-border focus:outline-none focus:border-blue-500"
            />
            <datalist id="recent-beneficiaries">
              {recentBeneficiaries.map((b) => (
                <option key={b} value={b} />
              ))}
            </datalist>
          </div>

          {/* Account Selection (Bank Account vs Cash) */}
          <div>
            <label className="text-xs font-semibold text-zinc-400 block mb-1">Payment Method</label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setAccount('Account')}
                className={`p-3 rounded-xl border text-xs font-semibold flex items-center justify-center gap-2 transition-all ${
                  account === 'Account'
                    ? 'bg-blue-600/20 border-blue-500/50 text-blue-300'
                    : 'bg-surface-elevated border-surface-border text-zinc-400 hover:text-white'
                }`}
              >
                <Building2 className="w-4 h-4 text-blue-400" />
                Account (Bank / MFS)
              </button>

              <button
                type="button"
                onClick={() => setAccount('Cash')}
                className={`p-3 rounded-xl border text-xs font-semibold flex items-center justify-center gap-2 transition-all ${
                  account === 'Cash'
                    ? 'bg-emerald-600/20 border-emerald-500/50 text-emerald-300'
                    : 'bg-surface-elevated border-surface-border text-zinc-400 hover:text-white'
                }`}
              >
                <Banknote className="w-4 h-4 text-emerald-400" />
                Cash
              </button>
            </div>
          </div>

          {/* Date Picker */}
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="text-xs font-semibold text-zinc-400">Date</label>
              <div className="flex items-center gap-1 text-[11px]">
                <button
                  type="button"
                  onClick={() => setDate(getTodayDateString())}
                  className="px-2 py-0.5 rounded bg-surface-elevated text-zinc-300 border border-surface-border hover:text-white"
                >
                  Today
                </button>
                <button
                  type="button"
                  onClick={setYesterday}
                  className="px-2 py-0.5 rounded bg-surface-elevated text-zinc-300 border border-surface-border hover:text-white"
                >
                  Yesterday
                </button>
              </div>
            </div>
            <div className="relative">
              <input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="w-full bg-surface-elevated rounded-xl p-3 text-sm text-white border border-surface-border focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>

          {/* Optional Notes */}
          <div>
            <label className="text-xs font-semibold text-zinc-400 block mb-1">
              Notes (Optional)
            </label>
            <input
              type="text"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Additional details..."
              className="w-full bg-surface-elevated rounded-xl p-2.5 text-xs text-white border border-surface-border focus:outline-none focus:border-blue-500"
            />
          </div>

          {/* Action Button */}
          <button
            type="submit"
            className={`w-full py-3.5 rounded-2xl font-bold text-sm tracking-wide shadow-lg active:scale-[0.98] transition-all ${
              type === 'income'
                ? 'bg-emerald-600 hover:bg-emerald-500 text-white'
                : 'bg-rose-600 hover:bg-rose-500 text-white'
            }`}
          >
            {editingTransaction ? 'Save Changes' : type === 'income' ? 'Save Income' : 'Save Expense'}
          </button>
        </form>
      </div>
    </div>
  );
};
