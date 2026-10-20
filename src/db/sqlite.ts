import initSqlJs, { Database } from 'sql.js';
import { loadSQLiteFromIDB, saveSQLiteToIDB } from './idbStorage';
import { Transaction, AccountSummary, MonthlySummary, SyncStats } from '../types';

let dbInstance: Database | null = null;
let saveDebounceTimer: number | null = null;
let currentDeviceId: string = '';

export function getDeviceId(): string {
  return currentDeviceId;
}

function debounceSave() {
  if (saveDebounceTimer) {
    clearTimeout(saveDebounceTimer);
  }
  saveDebounceTimer = window.setTimeout(() => {
    if (dbInstance) {
      const binary = dbInstance.export();
      saveSQLiteToIDB(binary);
    }
  }, 250);
}

export async function initDatabase(): Promise<Database> {
  if (dbInstance) return dbInstance;

  const initFn = typeof initSqlJs === 'function' ? initSqlJs : (initSqlJs as unknown as { default: typeof initSqlJs }).default;
  const SQL = await initFn({
    locateFile: (file) => {
      const baseUrl = import.meta.env.BASE_URL || './';
      return `${baseUrl.replace(/\/?$/, '/')}${file}`;
    }
  });

  const savedBinary = await loadSQLiteFromIDB();
  if (savedBinary && savedBinary.length > 0) {
    dbInstance = new SQL.Database(savedBinary);
  } else {
    dbInstance = new SQL.Database();
  }

  // Ensure tables exist
  dbInstance.run(`
    CREATE TABLE IF NOT EXISTS transactions (
      id TEXT PRIMARY KEY,
      date TEXT NOT NULL,
      type TEXT NOT NULL,
      benf TEXT NOT NULL,
      exp_cat TEXT NOT NULL,
      amount REAL NOT NULL,
      account TEXT NOT NULL,
      notes TEXT,
      created_at INTEGER NOT NULL,
      updated_at INTEGER NOT NULL,
      deleted_at INTEGER DEFAULT NULL,
      device_id TEXT NOT NULL,
      version INTEGER DEFAULT 1 NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_tx_date ON transactions(date);
    CREATE INDEX IF NOT EXISTS idx_tx_type ON transactions(type);
    CREATE INDEX IF NOT EXISTS idx_tx_cat ON transactions(exp_cat);
    CREATE INDEX IF NOT EXISTS idx_tx_deleted ON transactions(deleted_at);

    CREATE TABLE IF NOT EXISTS meta (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL
    );
  `);

  // Device ID initialization
  const metaRes = dbInstance.exec("SELECT value FROM meta WHERE key = 'device_id'");
  if (metaRes.length > 0 && metaRes[0].values.length > 0) {
    currentDeviceId = String(metaRes[0].values[0][0]);
  } else {
    currentDeviceId = crypto.randomUUID();
    dbInstance.run("INSERT INTO meta (key, value) VALUES ('device_id', ?)", [currentDeviceId]);
    debounceSave();
  }

  return dbInstance;
}

export function getAllTransactions(options?: {
  month?: string; // YYYY-MM
  type?: string;  // 'income' | 'expense' | 'all'
  search?: string;
  account?: string;
}): Transaction[] {
  if (!dbInstance) return [];

  let query = "SELECT * FROM transactions WHERE deleted_at IS NULL";
  const params: (string | number)[] = [];

  if (options?.month) {
    query += " AND date LIKE ?";
    params.push(`${options.month}%`);
  }

  if (options?.type && options.type !== 'all') {
    query += " AND type = ?";
    params.push(options.type);
  }

  if (options?.account && options.account !== 'all') {
    query += " AND account = ?";
    params.push(options.account);
  }

  if (options?.search && options.search.trim()) {
    query += " AND (benf LIKE ? OR exp_cat LIKE ? OR notes LIKE ?)";
    const term = `%${options.search.trim()}%`;
    params.push(term, term, term);
  }

  query += " ORDER BY date DESC, created_at DESC";

  const stmt = dbInstance.prepare(query);
  stmt.bind(params);

  const results: Transaction[] = [];
  while (stmt.step()) {
    const row = stmt.getAsObject() as unknown as Transaction;
    results.push(row);
  }
  stmt.free();

  return results;
}

export function addTransaction(data: {
  date: string;
  type: 'income' | 'expense';
  benf: string;
  exp_cat: string;
  amount: number;
  account: 'Account' | 'Cash';
  notes?: string;
}): Transaction {
  if (!dbInstance) throw new Error("Database not initialized");

  const id = crypto.randomUUID();
  const now = Date.now();
  const deviceId = currentDeviceId || crypto.randomUUID();

  dbInstance.run(`
    INSERT INTO transactions (
      id, date, type, benf, exp_cat, amount, account, notes,
      created_at, updated_at, deleted_at, device_id, version
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, 1)
  `, [
    id,
    data.date,
    data.type,
    data.benf.trim(),
    data.exp_cat.trim(),
    data.amount,
    data.account,
    data.notes?.trim() || null,
    now,
    now,
    deviceId
  ]);

  debounceSave();

  return {
    id,
    ...data,
    notes: data.notes?.trim() || undefined,
    created_at: now,
    updated_at: now,
    deleted_at: null,
    device_id: deviceId,
    version: 1
  };
}

export function updateTransaction(
  id: string,
  updates: Partial<Omit<Transaction, 'id' | 'created_at' | 'device_id'>>
): void {
  if (!dbInstance) throw new Error("Database not initialized");

  const now = Date.now();
  const fields: string[] = ['updated_at = ?', 'version = version + 1'];
  const values: (string | number | null)[] = [now];

  if (updates.date !== undefined) {
    fields.push('date = ?');
    values.push(updates.date);
  }
  if (updates.type !== undefined) {
    fields.push('type = ?');
    values.push(updates.type);
  }
  if (updates.benf !== undefined) {
    fields.push('benf = ?');
    values.push(updates.benf.trim());
  }
  if (updates.exp_cat !== undefined) {
    fields.push('exp_cat = ?');
    values.push(updates.exp_cat.trim());
  }
  if (updates.amount !== undefined) {
    fields.push('amount = ?');
    values.push(updates.amount);
  }
  if (updates.account !== undefined) {
    fields.push('account = ?');
    values.push(updates.account);
  }
  if (updates.notes !== undefined) {
    fields.push('notes = ?');
    values.push(updates.notes?.trim() || null);
  }

  values.push(id);

  dbInstance.run(`UPDATE transactions SET ${fields.join(', ')} WHERE id = ?`, values);
  debounceSave();
}

export function softDeleteTransaction(id: string): void {
  if (!dbInstance) throw new Error("Database not initialized");
  const now = Date.now();
  dbInstance.run(
    "UPDATE transactions SET deleted_at = ?, updated_at = ?, version = version + 1 WHERE id = ?",
    [now, now, id]
  );
  debounceSave();
}

export function getAccountSummary(month?: string): AccountSummary {
  if (!dbInstance) {
    return {
      totalIncome: 0,
      totalExpense: 0,
      totalBalance: 0,
      bankBalance: 0,
      cashBalance: 0,
      needBalancing: 0,
      isBalanced: true
    };
  }

  let filter = "WHERE deleted_at IS NULL";
  const params: string[] = [];
  if (month) {
    filter += " AND date LIKE ?";
    params.push(`${month}%`);
  }

  // Totals
  const res = dbInstance.exec(`
    SELECT
      SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as inc,
      SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as exp,
      SUM(CASE WHEN account = 'Account' THEN (CASE WHEN type = 'income' THEN amount ELSE -amount END) ELSE 0 END) as bank,
      SUM(CASE WHEN account = 'Cash' THEN (CASE WHEN type = 'income' THEN amount ELSE -amount END) ELSE 0 END) as cash
    FROM transactions ${filter}
  `, params);

  if (res.length === 0 || res[0].values.length === 0) {
    return {
      totalIncome: 0,
      totalExpense: 0,
      totalBalance: 0,
      bankBalance: 0,
      cashBalance: 0,
      needBalancing: 0,
      isBalanced: true
    };
  }

  const [inc, exp, bank, cash] = res[0].values[0];
  const totalIncome = Number(inc || 0);
  const totalExpense = Number(exp || 0);
  const totalBalance = totalIncome - totalExpense;
  const bankBalance = Number(bank || 0);
  const cashBalance = Number(cash || 0);
  const calculatedSum = bankBalance + cashBalance;
  const needBalancing = Math.round((totalBalance - calculatedSum) * 100) / 100;

  return {
    totalIncome,
    totalExpense,
    totalBalance,
    bankBalance,
    cashBalance,
    needBalancing,
    isBalanced: Math.abs(needBalancing) < 0.01
  };
}

export function getMonthlySummary(month: string): MonthlySummary {
  if (!dbInstance) {
    return { month, income: 0, expense: 0, balance: 0, categoryBreakdown: [] };
  }

  const summary = getAccountSummary(month);

  // Category breakdown for expenses
  const catRes = dbInstance.exec(`
    SELECT exp_cat, SUM(amount) as cat_total
    FROM transactions
    WHERE deleted_at IS NULL AND type = 'expense' AND date LIKE ?
    GROUP BY exp_cat
    ORDER BY cat_total DESC
  `, [`${month}%`]);

  const categoryBreakdown: { category: string; amount: number; percentage: number }[] = [];
  if (catRes.length > 0) {
    for (const row of catRes[0].values) {
      const category = String(row[0]);
      const amount = Number(row[1]);
      const percentage = summary.totalExpense > 0 ? (amount / summary.totalExpense) * 100 : 0;
      categoryBreakdown.push({
        category,
        amount,
        percentage: Math.round(percentage * 10) / 10
      });
    }
  }

  return {
    month,
    income: summary.totalIncome,
    expense: summary.totalExpense,
    balance: summary.totalBalance,
    categoryBreakdown
  };
}

export function getAvailableMonths(): string[] {
  if (!dbInstance) return [];

  const res = dbInstance.exec(`
    SELECT DISTINCT SUBSTR(date, 1, 7) as ym
    FROM transactions
    WHERE deleted_at IS NULL
    ORDER BY ym DESC
  `);

  if (res.length === 0) return [];
  return res[0].values.map(v => String(v[0]));
}

export function getRecentCategories(): string[] {
  if (!dbInstance) return [];
  const res = dbInstance.exec(`
    SELECT exp_cat, COUNT(*) as cnt
    FROM transactions
    WHERE deleted_at IS NULL
    GROUP BY exp_cat
    ORDER BY cnt DESC
    LIMIT 15
  `);
  if (res.length === 0) return [];
  return res[0].values.map(v => String(v[0]));
}

export function getRecentBeneficiaries(): string[] {
  if (!dbInstance) return [];
  const res = dbInstance.exec(`
    SELECT benf, COUNT(*) as cnt
    FROM transactions
    WHERE deleted_at IS NULL
    GROUP BY benf
    ORDER BY cnt DESC
    LIMIT 15
  `);
  if (res.length === 0) return [];
  return res[0].values.map(v => String(v[0]));
}

// -------------------------------------------------------------
// SYNC & BACKUP ENGINE (Last-Write-Wins with UUID)
// -------------------------------------------------------------

export function exportDatabaseBinary(): Uint8Array {
  if (!dbInstance) throw new Error("Database not initialized");
  return dbInstance.export();
}

export async function importDatabaseBinary(buffer: Uint8Array): Promise<void> {
  const initFn = typeof initSqlJs === 'function' ? initSqlJs : (initSqlJs as unknown as { default: typeof initSqlJs }).default;
  const SQL = await initFn({
    locateFile: (file) => {
      const baseUrl = import.meta.env.BASE_URL || './';
      return `${baseUrl.replace(/\/?$/, '/')}${file}`;
    }
  });
  dbInstance = new SQL.Database(buffer);
  await saveSQLiteToIDB(buffer);
}

export function exportSyncPackage(): {
  deviceId: string;
  timestamp: number;
  records: Transaction[];
} {
  if (!dbInstance) throw new Error("Database not initialized");
  const stmt = dbInstance.prepare("SELECT * FROM transactions");
  const records: Transaction[] = [];
  while (stmt.step()) {
    records.push(stmt.getAsObject() as unknown as Transaction);
  }
  stmt.free();

  return {
    deviceId: currentDeviceId,
    timestamp: Date.now(),
    records
  };
}

export function importAndMergeSyncPackage(pkg: {
  deviceId: string;
  timestamp: number;
  records: Transaction[];
}): { inserted: number; updated: number; skipped: number } {
  if (!dbInstance) throw new Error("Database not initialized");

  let inserted = 0;
  let updated = 0;
  let skipped = 0;

  for (const remote of pkg.records) {
    const res = dbInstance.exec("SELECT updated_at, deleted_at, version FROM transactions WHERE id = ?", [remote.id]);

    if (res.length === 0 || res[0].values.length === 0) {
      // Record does not exist locally -> INSERT
      dbInstance.run(`
        INSERT INTO transactions (
          id, date, type, benf, exp_cat, amount, account, notes,
          created_at, updated_at, deleted_at, device_id, version
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
      `, [
        remote.id,
        remote.date,
        remote.type,
        remote.benf,
        remote.exp_cat,
        remote.amount,
        remote.account,
        remote.notes || null,
        remote.created_at,
        remote.updated_at,
        remote.deleted_at ?? null,
        remote.device_id,
        remote.version
      ]);
      inserted++;
    } else {
      // Record exists locally -> Compare timestamps (Last-Write-Wins)
      const localUpdatedAt = Number(res[0].values[0][0]);
      if (remote.updated_at > localUpdatedAt) {
        dbInstance.run(`
          UPDATE transactions SET
            date = ?,
            type = ?,
            benf = ?,
            exp_cat = ?,
            amount = ?,
            account = ?,
            notes = ?,
            updated_at = ?,
            deleted_at = ?,
            device_id = ?,
            version = ?
          WHERE id = ?
        `, [
          remote.date,
          remote.type,
          remote.benf,
          remote.exp_cat,
          remote.amount,
          remote.account,
          remote.notes || null,
          remote.updated_at,
          remote.deleted_at ?? null,
          remote.device_id,
          remote.version,
          remote.id
        ]);
        updated++;
      } else {
        skipped++;
      }
    }
  }

  debounceSave();
  return { inserted, updated, skipped };
}

export function exportCSV(): string {
  if (!dbInstance) return '';
  const txs = getAllTransactions();
  const headers = ['id', 'date', 'type', 'beneficiary', 'category', 'amount', 'account', 'notes'];
  const rows = txs.map(t => [
    t.id,
    t.date,
    t.type,
    `"${(t.benf || '').replace(/"/g, '""')}"`,
    `"${(t.exp_cat || '').replace(/"/g, '""')}"`,
    t.amount.toFixed(2),
    t.account,
    `"${(t.notes || '').replace(/"/g, '""')}"`
  ]);

  return [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
}

export function clearAllData(): void {
  if (!dbInstance) throw new Error("Database not initialized");
  dbInstance.run("DELETE FROM transactions");
  debounceSave();
}

export function getSyncStats(): SyncStats {
  if (!dbInstance) {
    return { deviceId: currentDeviceId, totalRecords: 0, activeRecords: 0, deletedRecords: 0, lastUpdated: 0 };
  }
  const totalRes = dbInstance.exec("SELECT COUNT(*) FROM transactions");
  const activeRes = dbInstance.exec("SELECT COUNT(*) FROM transactions WHERE deleted_at IS NULL");
  const delRes = dbInstance.exec("SELECT COUNT(*) FROM transactions WHERE deleted_at IS NOT NULL");
  const latestRes = dbInstance.exec("SELECT MAX(updated_at) FROM transactions");

  return {
    deviceId: currentDeviceId,
    totalRecords: Number(totalRes[0]?.values[0]?.[0] || 0),
    activeRecords: Number(activeRes[0]?.values[0]?.[0] || 0),
    deletedRecords: Number(delRes[0]?.values[0]?.[0] || 0),
    lastUpdated: Number(latestRes[0]?.values[0]?.[0] || 0)
  };
}
