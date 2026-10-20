import { useEffect, useState, useCallback } from 'react';
import {
  initDatabase,
  getAllTransactions,
  getAccountSummary,
  getMonthlySummary,
  getAvailableMonths,
  addTransaction,
  updateTransaction,
  softDeleteTransaction,
  getRecentCategories,
  getRecentBeneficiaries
} from './db/sqlite';
import { Transaction, AccountSummary, MonthlySummary, TransactionType, AccountType } from './types';
import { getCurrentMonthString, formatMonthDisplay } from './utils/formatters';
import { Header } from './components/Header';
import { MonthSelector } from './components/MonthSelector';
import { BalanceCard } from './components/BalanceCard';
import { TransactionList } from './components/TransactionList';
import { AddTransactionSheet } from './components/AddTransactionSheet';
import { AnalyticsView } from './components/AnalyticsView';
import { SettingsView } from './components/SettingsView';
import { BottomNav, NavTab } from './components/BottomNav';
import { Search, Loader2, Wallet } from 'lucide-react';

export function App() {
  const [dbReady, setDbReady] = useState<boolean>(false);
  const [initError, setInitError] = useState<string | null>(null);

  // Navigation and Modal State
  const [currentTab, setCurrentTab] = useState<NavTab>('dashboard');
  const [isAddOpen, setIsAddOpen] = useState<boolean>(false);
  const [editingTx, setEditingTx] = useState<Transaction | null>(null);

  // Filter State
  const [selectedMonth, setSelectedMonth] = useState<string>(getCurrentMonthString());
  const [availableMonths, setAvailableMonths] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [typeFilter, setTypeFilter] = useState<'all' | 'income' | 'expense'>('all');
  const [accountFilter, setAccountFilter] = useState<'all' | 'Account' | 'Cash'>('all');

  // Data State
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accountSummary, setAccountSummary] = useState<AccountSummary>({
    totalIncome: 0,
    totalExpense: 0,
    totalBalance: 0,
    bankBalance: 0,
    cashBalance: 0,
    needBalancing: 0,
    isBalanced: true
  });
  const [monthlySummary, setMonthlySummary] = useState<MonthlySummary>({
    month: getCurrentMonthString(),
    income: 0,
    expense: 0,
    balance: 0,
    categoryBreakdown: []
  });
  const [recentCats, setRecentCats] = useState<string[]>([]);
  const [recentBenfs, setRecentBenfs] = useState<string[]>([]);

  // 1. Initialize SQLite Database
  useEffect(() => {
    initDatabase()
      .then(() => {
        const months = getAvailableMonths();
        setAvailableMonths(months);
        if (months.length > 0 && !months.includes(selectedMonth)) {
          setSelectedMonth(months[0]);
        }
        setDbReady(true);
      })
      .catch((err) => {
        console.error('SQLite init error:', err);
        setInitError(err.message || String(err));
      });
  }, []);

  // 2. Fetch / Refresh Data
  const refreshData = useCallback(() => {
    if (!dbReady) return;

    const months = getAvailableMonths();
    setAvailableMonths(months);

    const monthParam = selectedMonth === 'all' ? undefined : selectedMonth;
    const txs = getAllTransactions({
      month: monthParam,
      type: typeFilter,
      account: accountFilter,
      search: searchQuery
    });
    setTransactions(txs);

    const accSummary = getAccountSummary(monthParam);
    setAccountSummary(accSummary);

    const monSummary = getMonthlySummary(monthParam || months[0] || getCurrentMonthString());
    setMonthlySummary(monSummary);

    setRecentCats(getRecentCategories());
    setRecentBenfs(getRecentBeneficiaries());
  }, [dbReady, selectedMonth, typeFilter, accountFilter, searchQuery]);

  useEffect(() => {
    refreshData();
  }, [refreshData]);

  // Handlers
  const handleAddOrEdit = (data: {
    id?: string;
    date: string;
    type: TransactionType;
    benf: string;
    exp_cat: string;
    amount: number;
    account: AccountType;
    notes?: string;
  }) => {
    if (data.id) {
      updateTransaction(data.id, data);
    } else {
      addTransaction(data);
    }
    refreshData();
  };

  const handleDelete = (id: string) => {
    softDeleteTransaction(id);
    refreshData();
  };

  const handleEdit = (tx: Transaction) => {
    setEditingTx(tx);
    setIsAddOpen(true);
  };

  // Loading Screen
  if (!dbReady && !initError) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center p-6 text-center">
        <div className="w-16 h-16 rounded-3xl bg-blue-600/20 border border-blue-500/30 flex items-center justify-center text-blue-400 mb-4 animate-pulse">
          <Wallet className="w-8 h-8" />
        </div>
        <h2 className="text-lg font-bold text-white mb-2">Starting SQLite Engine...</h2>
        <p className="text-xs text-zinc-400 max-w-xs mb-6">
          Initializing in-browser WebAssembly SQLite and restoring local records.
        </p>
        <Loader2 className="w-6 h-6 text-blue-500 animate-spin" />
      </div>
    );
  }

  // Error Screen
  if (initError) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center p-6 text-center">
        <div className="text-rose-500 font-bold mb-2">Database Initialization Error</div>
        <p className="text-xs text-zinc-400 mb-4 max-w-sm">{initError}</p>
        <button
          onClick={() => window.location.reload()}
          className="px-4 py-2 bg-blue-600 rounded-xl text-xs font-semibold text-white"
        >
          Reload App
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background text-zinc-100 flex flex-col max-w-md mx-auto relative antialiased selection:bg-blue-500/30">
      {/* Top Header */}
      <Header
        onOpenSync={() => setCurrentTab('sync')}
        onRefresh={refreshData}
      />

      {/* Main Content Area */}
      <main className="flex-1 px-4 py-4 pb-28 space-y-4 overflow-y-auto">
        {/* Month Selector Carousel (Visible in Dashboard & Ledger & Analytics) */}
        {currentTab !== 'sync' && (
          <MonthSelector
            selectedMonth={selectedMonth}
            availableMonths={availableMonths}
            onChange={setSelectedMonth}
          />
        )}

        {/* TAB 1: DASHBOARD (HOME) */}
        {currentTab === 'dashboard' && (
          <>
            <BalanceCard
              summary={accountSummary}
              monthLabel={selectedMonth === 'all' ? 'All Time' : formatMonthDisplay(selectedMonth)}
            />

            <div className="flex items-center justify-between pt-2">
              <h2 className="text-sm font-bold text-white tracking-tight">Recent Transactions</h2>
              <button
                onClick={() => setCurrentTab('ledger')}
                className="text-xs font-semibold text-blue-400 hover:text-blue-300"
              >
                View all ({transactions.length})
              </button>
            </div>

            <TransactionList
              transactions={transactions.slice(0, 10)}
              onDelete={handleDelete}
              onEdit={handleEdit}
            />
          </>
        )}

        {/* TAB 2: LEDGER (FULL HISTORY & SEARCH) */}
        {currentTab === 'ledger' && (
          <div className="space-y-3">
            {/* Search Bar */}
            <div className="relative">
              <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-zinc-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search transactions..."
                className="w-full bg-surface rounded-2xl pl-10 pr-4 py-3 text-xs text-white placeholder-zinc-500 border border-surface-border focus:outline-none focus:border-blue-500"
              />
            </div>

            {/* Quick Filters */}
            <div className="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs">
              <button
                onClick={() => setTypeFilter('all')}
                className={`px-3 py-1.5 rounded-xl border flex-shrink-0 font-medium transition-all ${
                  typeFilter === 'all'
                    ? 'bg-zinc-200 text-zinc-900 border-zinc-200 font-bold'
                    : 'bg-surface text-zinc-400 border-surface-border hover:text-white'
                }`}
              >
                All Types
              </button>
              <button
                onClick={() => setTypeFilter('expense')}
                className={`px-3 py-1.5 rounded-xl border flex-shrink-0 font-medium transition-all ${
                  typeFilter === 'expense'
                    ? 'bg-rose-600 text-white border-rose-600 font-bold'
                    : 'bg-surface text-zinc-400 border-surface-border hover:text-white'
                }`}
              >
                Expenses
              </button>
              <button
                onClick={() => setTypeFilter('income')}
                className={`px-3 py-1.5 rounded-xl border flex-shrink-0 font-medium transition-all ${
                  typeFilter === 'income'
                    ? 'bg-emerald-600 text-white border-emerald-600 font-bold'
                    : 'bg-surface text-zinc-400 border-surface-border hover:text-white'
                }`}
              >
                Income
              </button>

              <span className="text-zinc-600">|</span>

              <button
                onClick={() => setAccountFilter(accountFilter === 'Account' ? 'all' : 'Account')}
                className={`px-3 py-1.5 rounded-xl border flex-shrink-0 font-medium transition-all ${
                  accountFilter === 'Account'
                    ? 'bg-blue-600 text-white border-blue-600 font-bold'
                    : 'bg-surface text-zinc-400 border-surface-border hover:text-white'
                }`}
              >
                Bank
              </button>
              <button
                onClick={() => setAccountFilter(accountFilter === 'Cash' ? 'all' : 'Cash')}
                className={`px-3 py-1.5 rounded-xl border flex-shrink-0 font-medium transition-all ${
                  accountFilter === 'Cash'
                    ? 'bg-emerald-600 text-white border-emerald-600 font-bold'
                    : 'bg-surface text-zinc-400 border-surface-border hover:text-white'
                }`}
              >
                Cash
              </button>
            </div>

            <TransactionList
              transactions={transactions}
              onDelete={handleDelete}
              onEdit={handleEdit}
            />
          </div>
        )}

        {/* TAB 3: ANALYTICS */}
        {currentTab === 'analytics' && (
          <AnalyticsView
            monthlySummary={monthlySummary}
            overallSummary={accountSummary}
            month={selectedMonth}
          />
        )}

        {/* TAB 4: SYNC & SETTINGS */}
        {currentTab === 'sync' && (
          <SettingsView onDataChanged={refreshData} />
        )}
      </main>

      {/* Floating / Sticky Bottom Navigation */}
      <BottomNav
        currentTab={currentTab}
        onTabChange={setCurrentTab}
        onOpenAdd={() => {
          setEditingTx(null);
          setIsAddOpen(true);
        }}
      />

      {/* Quick Add Bottom Sheet Modal */}
      <AddTransactionSheet
        isOpen={isAddOpen}
        onClose={() => {
          setIsAddOpen(false);
          setEditingTx(null);
        }}
        onSubmit={handleAddOrEdit}
        recentCategories={recentCats}
        recentBeneficiaries={recentBenfs}
        editingTransaction={editingTx}
      />
    </div>
  );
}

export default App;
