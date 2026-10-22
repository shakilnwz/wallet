package db

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/shakilnwz/wallet/core/models"
)

func openTestDB(t *testing.T, name string) *DB {
	t.Helper()
	d, err := Open(filepath.Join(t.TempDir(), name))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	return d
}

func addTx(t *testing.T, d *DB, date string, typ models.TransactionType, cat string, paisa int64, acc models.AccountType) *models.Transaction {
	t.Helper()
	tx, err := d.AddTransaction(models.Transaction{
		Date: date, Type: typ, Benf: cat, ExpCat: cat, Amount: paisa, Account: acc,
	})
	if err != nil {
		t.Fatalf("AddTransaction: %v", err)
	}
	return tx
}

func TestMigrateLegacyFloatAmounts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")

	// Build a database exactly as pre-versioning builds left it.
	raw, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = raw.Exec(`
	CREATE TABLE transactions (
		id TEXT PRIMARY KEY, date TEXT NOT NULL, type TEXT NOT NULL, benf TEXT NOT NULL,
		exp_cat TEXT NOT NULL, amount REAL NOT NULL, account TEXT NOT NULL, notes TEXT,
		created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, deleted_at INTEGER DEFAULT NULL,
		device_id TEXT NOT NULL, version INTEGER DEFAULT 1 NOT NULL
	);
	CREATE TABLE meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
	INSERT INTO meta VALUES ('device_id', 'legacy-device');
	INSERT INTO transactions VALUES
		('a', '2026-10-01', 'expense', 'x', 'Food', 0.29, 'Cash', NULL, 1, 1, NULL, 'legacy-device', 1),
		('b', '2026-10-02', 'income', 'y', 'Salary', 1234.56, 'Account', 'n', 1, 1, NULL, 'legacy-device', 1),
		('c', '2026-10-03', 'expense', 'z', 'Food', 0.1, 'Cash', NULL, 1, 1, 5, 'legacy-device', 2);
	`)
	if err != nil {
		t.Fatal(err)
	}
	_ = raw.Close()

	for i := 0; i < 2; i++ { // second open must be a no-op
		d, err := Open(path)
		if err != nil {
			t.Fatalf("Open #%d: %v", i+1, err)
		}
		if d.DeviceID() != "legacy-device" {
			t.Errorf("device id lost: %q", d.DeviceID())
		}
		pkg, err := d.ExportSyncPackage()
		if err != nil {
			t.Fatal(err)
		}
		want := map[string]int64{"a": 29, "b": 123456, "c": 10}
		for _, r := range pkg.Records {
			if r.Amount != want[r.ID] {
				t.Errorf("open #%d: %s amount = %d paisa, want %d", i+1, r.ID, r.Amount, want[r.ID])
			}
		}
		if len(pkg.Records) != 3 {
			t.Errorf("got %d records, want 3", len(pkg.Records))
		}
		var ver string
		if err := d.db.QueryRow("SELECT value FROM meta WHERE key = 'schema_version'").Scan(&ver); err != nil || ver != "2" {
			t.Errorf("schema_version = %q (%v), want 2", ver, err)
		}
		_ = d.Close()
	}
}

func TestOpenRejectsNewerSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "future.db")
	d, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = d.db.Exec("UPDATE meta SET value = '99' WHERE key = 'schema_version'")
	_ = d.Close()

	if d, err := Open(path); err == nil {
		_ = d.Close()
		t.Fatal("expected Open to refuse a newer schema")
	}
}

func TestPaisaSumsAreExact(t *testing.T) {
	d := openTestDB(t, "exact.db")
	addTx(t, d, "2026-10-01", models.TypeExpense, "Food", 10, models.AccountCash)
	addTx(t, d, "2026-10-01", models.TypeExpense, "Food", 20, models.AccountCash)

	s, err := d.GetAccountSummary("")
	if err != nil {
		t.Fatal(err)
	}
	if s.TotalExpense != 30 || s.CashBalance != -30 || !s.IsBalanced {
		t.Errorf("got expense=%d cash=%d balanced=%v; want 30, -30, true", s.TotalExpense, s.CashBalance, s.IsBalanced)
	}
}

func TestMonthlySummaryAllTime(t *testing.T) {
	d := openTestDB(t, "all.db")
	addTx(t, d, "2026-09-15", models.TypeExpense, "Rent", 1000_00, models.AccountBank)
	addTx(t, d, "2026-10-01", models.TypeExpense, "Food", 300_00, models.AccountCash)
	addTx(t, d, "2026-10-02", models.TypeExpense, "Food", 200_00, models.AccountCash)

	for _, month := range []string{"all", ""} {
		m, err := d.GetMonthlySummary(month)
		if err != nil {
			t.Fatal(err)
		}
		if m.Expense != 1500_00 || len(m.Categories) != 2 {
			t.Fatalf("month %q: expense=%d categories=%d; want 150000, 2", month, m.Expense, len(m.Categories))
		}
		if m.Categories[0].Category != "Rent" || m.Categories[0].Percentage != 66.7 {
			t.Errorf("month %q: top category = %+v", month, m.Categories[0])
		}
	}

	oct, err := d.GetMonthlySummary("2026-10")
	if err != nil {
		t.Fatal(err)
	}
	if oct.Expense != 500_00 || len(oct.Categories) != 1 {
		t.Errorf("2026-10: expense=%d categories=%d; want 50000, 1", oct.Expense, len(oct.Categories))
	}
}

func TestClearAllDataPropagatesToPeer(t *testing.T) {
	local := openTestDB(t, "local.db")
	peer := openTestDB(t, "peer.db")

	addTx(t, local, "2026-10-01", models.TypeExpense, "Food", 100_00, models.AccountCash)
	addTx(t, local, "2026-10-02", models.TypeIncome, "Salary", 500_00, models.AccountBank)

	pkg, err := local.ExportSyncPackage()
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := peer.ImportAndMergeSyncPackage(pkg); err != nil {
		t.Fatal(err)
	}

	n, err := local.ClearAllData()
	if err != nil || n != 2 {
		t.Fatalf("ClearAllData = %d, %v; want 2", n, err)
	}
	if again, _ := local.ClearAllData(); again != 0 {
		t.Errorf("second ClearAllData deleted %d rows, want 0", again)
	}

	// Peer pulls the tombstones: its copies must be deleted, not resurrected.
	pkg, err = local.ExportSyncPackage()
	if err != nil {
		t.Fatal(err)
	}
	if _, upd, _, err := peer.ImportAndMergeSyncPackage(pkg); err != nil || upd != 2 {
		t.Fatalf("merge tombstones: updated=%d err=%v; want 2", upd, err)
	}
	for name, d := range map[string]*DB{"local": local, "peer": peer} {
		txs, err := d.GetTransactions(QueryFilter{})
		if err != nil {
			t.Fatal(err)
		}
		if len(txs) != 0 {
			t.Errorf("%s still has %d live transactions", name, len(txs))
		}
		if st := d.GetSyncStats(); st.DeletedRecords != 2 {
			t.Errorf("%s: %d tombstones, want 2 (rows must be kept)", name, st.DeletedRecords)
		}
	}
}

func TestPragmasApplyToEveryConnection(t *testing.T) {
	d := openTestDB(t, "pragma.db")
	ctx := context.Background()

	// Hold two connections at once so the pool can't hand back the same one.
	var conns []*sql.Conn
	for i := 0; i < 2; i++ {
		c, err := d.db.Conn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		conns = append(conns, c)
	}
	for i, c := range conns {
		var timeout int
		var mode string
		if err := c.QueryRowContext(ctx, "PRAGMA busy_timeout").Scan(&timeout); err != nil {
			t.Fatal(err)
		}
		if err := c.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&mode); err != nil {
			t.Fatal(err)
		}
		if timeout != 5000 || mode != "wal" {
			t.Errorf("conn %d: busy_timeout=%d journal_mode=%q; want 5000, wal", i, timeout, mode)
		}
	}
}
