export type TransactionType = 'income' | 'expense';
export type AccountType = 'Account' | 'Cash';

export interface Transaction {
  id: string;               // UUID v4
  date: string;             // YYYY-MM-DD
  type: TransactionType;    // 'income' | 'expense'
  benf: string;             // Beneficiary / Payee / Source
  exp_cat: string;          // Category (Rent, Meal, Salary...)
  amount: number;           // Amount in BDT (৳)
  account: AccountType;     // 'Account' (Bank/MFS) | 'Cash'
  notes?: string;           // Optional note/remarks
  created_at: number;       // Unix epoch ms
  updated_at: number;       // Unix epoch ms
  deleted_at: number | null;// Tombstone timestamp for syncing
  device_id: string;        // Originating device ID
  version: number;          // Monotonic version counter
}

export interface AccountSummary {
  totalIncome: number;
  totalExpense: number;
  totalBalance: number;
  bankBalance: number;
  cashBalance: number;
  needBalancing: number;
  isBalanced: boolean;
}

export interface MonthlySummary {
  month: string;            // YYYY-MM
  income: number;
  expense: number;
  balance: number;
  categoryBreakdown: { category: string; amount: number; percentage: number }[];
}

export interface SyncStats {
  deviceId: string;
  totalRecords: number;
  activeRecords: number;
  deletedRecords: number;
  lastUpdated: number;
}
