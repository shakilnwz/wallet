package db

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shakilnwz/wallet/core/models"
	_ "modernc.org/sqlite"
)

type DB struct {
	db       *sql.DB
	deviceID string
	mu       sync.RWMutex
}

func Open(customPath string) (*DB, error) {
	dbPath := customPath
	if dbPath == "" {
		// On Android, Gio points XDG_CONFIG_HOME at the app's private files dir,
		// so this resolves to <filesDir>/config/wallet/wallet.db there.
		baseDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to locate config dir: %w", err)
		}
		appDir := filepath.Join(baseDir, "wallet")
		if err := os.MkdirAll(appDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create data dir: %w", err)
		}
		dbPath = filepath.Join(appDir, "wallet.db")
	}

	// Pragmas go in the DSN so they apply to every pooled connection, not just
	// the first one. _txlock=immediate takes the write lock at BEGIN, avoiding
	// lock-upgrade deadlocks between concurrent transactions.
	dsn := dbPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_txlock=immediate"
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	var journalMode string
	if err := sqlDB.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}
	if journalMode != "wal" {
		log.Printf("warning: sqlite journal_mode is %q, expected wal", journalMode)
	}

	d := &DB{db: sqlDB}
	if err := d.migrate(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	if err := d.initDeviceID(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to init device ID: %w", err)
	}

	return d, nil
}

func (d *DB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.db.Close()
}

func (d *DB) DeviceID() string {
	return d.deviceID
}

// migrations[i] upgrades the schema from version i to i+1. Append only; never
// edit a migration that has shipped.
var migrations = []func(*sql.Tx) error{
	migrateBaseSchema,
	migrateAmountToPaisa,
}

func (d *DB) migrate() error {
	if _, err := d.db.Exec("CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL)"); err != nil {
		return err
	}

	var current int
	var raw string
	err := d.db.QueryRow("SELECT value FROM meta WHERE key = 'schema_version'").Scan(&raw)
	switch {
	case err == sql.ErrNoRows:
		current = 0
	case err != nil:
		return err
	default:
		if current, err = strconv.Atoi(raw); err != nil {
			return fmt.Errorf("bad schema_version %q: %w", raw, err)
		}
	}
	if current > len(migrations) {
		return fmt.Errorf("database schema v%d is newer than this app supports (v%d)", current, len(migrations))
	}

	for v := current; v < len(migrations); v++ {
		tx, err := d.db.Begin()
		if err != nil {
			return err
		}
		if err := migrations[v](tx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration to v%d: %w", v+1, err)
		}
		if _, err := tx.Exec("INSERT OR REPLACE INTO meta (key, value) VALUES ('schema_version', ?)", strconv.Itoa(v+1)); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// migrateBaseSchema is the original schema. IF NOT EXISTS keeps it a no-op on
// databases created before schema versioning existed.
func migrateBaseSchema(tx *sql.Tx) error {
	schema := `
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
	`
	_, err := tx.Exec(schema)
	return err
}

// migrateAmountToPaisa replaces the REAL amount column with exact integer paisa.
func migrateAmountToPaisa(tx *sql.Tx) error {
	stmts := []string{
		"ALTER TABLE transactions ADD COLUMN amount_paisa INTEGER NOT NULL DEFAULT 0",
		"UPDATE transactions SET amount_paisa = CAST(ROUND(amount * 100) AS INTEGER)",
		"ALTER TABLE transactions DROP COLUMN amount",
	}
	for _, s := range stmts {
		if _, err := tx.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// monthClause returns a SQL condition (with leading AND) restricting rows to a
// YYYY-MM month, or nothing when month is empty or "all".
func monthClause(month string) (string, []any) {
	if month == "" || month == "all" {
		return "", nil
	}
	return " AND date LIKE ?", []any{month + "%"}
}

func (d *DB) initDeviceID() error {
	var val string
	err := d.db.QueryRow("SELECT value FROM meta WHERE key = 'device_id'").Scan(&val)
	if err == sql.ErrNoRows {
		d.deviceID = uuid.NewString()
		_, err = d.db.Exec("INSERT INTO meta (key, value) VALUES ('device_id', ?)", d.deviceID)
		return err
	}
	if err != nil {
		return err
	}
	d.deviceID = val
	return nil
}

// -------------------------------------------------------------
// Transaction CRUD
// -------------------------------------------------------------

type QueryFilter struct {
	Month   string // YYYY-MM or empty for all
	Type    string // income | expense | all
	Account string // Account | Cash | all
	Search  string // text search
}

func (d *DB) GetTransactions(filter QueryFilter) ([]models.Transaction, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("SELECT " + txColumns + " ")
	sb.WriteString("FROM transactions WHERE deleted_at IS NULL")

	clause, args := monthClause(filter.Month)
	sb.WriteString(clause + " ")

	if filter.Type != "" && filter.Type != "all" {
		sb.WriteString("AND type = ? ")
		args = append(args, filter.Type)
	}

	if filter.Account != "" && filter.Account != "all" {
		sb.WriteString("AND account = ? ")
		args = append(args, filter.Account)
	}

	if filter.Search != "" {
		term := "%" + strings.TrimSpace(filter.Search) + "%"
		sb.WriteString("AND (benf LIKE ? OR exp_cat LIKE ? OR notes LIKE ?) ")
		args = append(args, term, term, term)
	}

	sb.WriteString("ORDER BY date DESC, created_at DESC")

	rows, err := d.db.Query(sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTransactions(rows)
}

const txColumns = "id, date, type, benf, exp_cat, amount_paisa, account, COALESCE(notes, ''), created_at, updated_at, deleted_at, device_id, version"

// scanTransactions reads rows selected with txColumns.
func scanTransactions(rows *sql.Rows) ([]models.Transaction, error) {
	var results []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var delAt sql.NullInt64
		err := rows.Scan(
			&tx.ID, &tx.Date, &tx.Type, &tx.Benf, &tx.ExpCat, &tx.Amount,
			&tx.Account, &tx.Notes, &tx.CreatedAt, &tx.UpdatedAt, &delAt,
			&tx.DeviceID, &tx.Version,
		)
		if err != nil {
			return nil, err
		}
		if delAt.Valid {
			tx.DeletedAt = &delAt.Int64
		}
		results = append(results, tx)
	}
	return results, rows.Err()
}

func (d *DB) AddTransaction(tx models.Transaction) (*models.Transaction, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UnixMilli()
	if tx.ID == "" {
		tx.ID = uuid.NewString()
	}
	tx.CreatedAt = now
	tx.UpdatedAt = now
	tx.DeletedAt = nil
	tx.DeviceID = d.deviceID
	tx.Version = 1

	query := `
	INSERT INTO transactions (
		id, date, type, benf, exp_cat, amount_paisa, account, notes,
		created_at, updated_at, deleted_at, device_id, version
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, 1)
	`

	_, err := d.db.Exec(query,
		tx.ID, tx.Date, string(tx.Type), strings.TrimSpace(tx.Benf),
		strings.TrimSpace(tx.ExpCat), tx.Amount, string(tx.Account),
		strings.TrimSpace(tx.Notes), tx.CreatedAt, tx.UpdatedAt, tx.DeviceID,
	)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (d *DB) UpdateTransaction(tx models.Transaction) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UnixMilli()
	query := `
	UPDATE transactions SET
		date = ?,
		type = ?,
		benf = ?,
		exp_cat = ?,
		amount_paisa = ?,
		account = ?,
		notes = ?,
		updated_at = ?,
		version = version + 1
	WHERE id = ? AND deleted_at IS NULL
	`

	_, err := d.db.Exec(query,
		tx.Date, string(tx.Type), strings.TrimSpace(tx.Benf),
		strings.TrimSpace(tx.ExpCat), tx.Amount, string(tx.Account),
		strings.TrimSpace(tx.Notes), now, tx.ID,
	)
	return err
}

func (d *DB) SoftDeleteTransaction(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UnixMilli()
	query := `
	UPDATE transactions SET
		deleted_at = ?,
		updated_at = ?,
		version = version + 1
	WHERE id = ?
	`
	_, err := d.db.Exec(query, now, now, id)
	return err
}

// ClearAllData tombstones every live transaction and returns how many were
// deleted. Tombstones sync like any other edit, so peers delete them too on
// their next sync. Rows are kept (soft-deleted) and remain recoverable.
func (d *DB) ClearAllData() (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now().UnixMilli()
	res, err := d.db.Exec(`
	UPDATE transactions SET
		deleted_at = ?,
		updated_at = ?,
		version = version + 1
	WHERE deleted_at IS NULL
	`, now, now)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// -------------------------------------------------------------
// Financial Summaries & Calculations
// -------------------------------------------------------------

func (d *DB) GetAccountSummary(month string) (*models.AccountSummary, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString(`
	SELECT
		COALESCE(SUM(CASE WHEN type = 'income' THEN amount_paisa ELSE 0 END), 0) as inc,
		COALESCE(SUM(CASE WHEN type = 'expense' THEN amount_paisa ELSE 0 END), 0) as exp,
		COALESCE(SUM(CASE WHEN account = 'Account' THEN (CASE WHEN type = 'income' THEN amount_paisa ELSE -amount_paisa END) ELSE 0 END), 0) as bank,
		COALESCE(SUM(CASE WHEN account = 'Cash' THEN (CASE WHEN type = 'income' THEN amount_paisa ELSE -amount_paisa END) ELSE 0 END), 0) as cash
	FROM transactions
	WHERE deleted_at IS NULL
	`)

	clause, args := monthClause(month)
	sb.WriteString(clause)

	var inc, exp, bank, cash int64
	err := d.db.QueryRow(sb.String(), args...).Scan(&inc, &exp, &bank, &cash)
	if err != nil {
		return nil, err
	}

	totalBal := inc - exp
	diff := totalBal - (bank + cash)

	return &models.AccountSummary{
		TotalIncome:   inc,
		TotalExpense:  exp,
		TotalBalance:  totalBal,
		BankBalance:   bank,
		CashBalance:   cash,
		NeedBalancing: diff,
		IsBalanced:    diff == 0,
	}, nil
}

func (d *DB) GetMonthlySummary(month string) (*models.MonthlySummary, error) {
	accSummary, err := d.GetAccountSummary(month)
	if err != nil {
		return nil, err
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	clause, args := monthClause(month)
	query := `
	SELECT exp_cat, SUM(amount_paisa) as total
	FROM transactions
	WHERE deleted_at IS NULL AND type = 'expense'` + clause + `
	GROUP BY exp_cat
	ORDER BY total DESC
	`

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryBreakdown
	for rows.Next() {
		var cat string
		var amount int64
		if err := rows.Scan(&cat, &amount); err != nil {
			return nil, err
		}
		pct := 0.0
		if accSummary.TotalExpense > 0 {
			pct = math.Round(float64(amount)/float64(accSummary.TotalExpense)*1000) / 10
		}
		categories = append(categories, models.CategoryBreakdown{
			Category:   cat,
			Amount:     amount,
			Percentage: pct,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.MonthlySummary{
		Month:      month,
		Income:     accSummary.TotalIncome,
		Expense:    accSummary.TotalExpense,
		Balance:    accSummary.TotalBalance,
		Categories: categories,
	}, nil
}

func (d *DB) GetAvailableMonths() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	query := `
	SELECT DISTINCT SUBSTR(date, 1, 7) as ym
	FROM transactions
	WHERE deleted_at IS NULL
	ORDER BY ym DESC
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanStrings(rows)
}

func (d *DB) GetRecentCategories() ([]string, error) {
	return d.mostUsed("exp_cat")
}

func (d *DB) GetRecentBeneficiaries() ([]string, error) {
	return d.mostUsed("benf")
}

// mostUsed returns the 15 most frequent values of column among live rows.
// column must be a trusted identifier, never user input.
func (d *DB) mostUsed(column string) ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT ` + column + `
		FROM transactions
		WHERE deleted_at IS NULL
		GROUP BY ` + column + `
		ORDER BY COUNT(*) DESC
		LIMIT 15
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanStrings(rows)
}

func scanStrings(rows *sql.Rows) ([]string, error) {
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// -------------------------------------------------------------
// Sync & Backup Engine (UUID + LWW)
// -------------------------------------------------------------

func (d *DB) GetSyncStats() models.SyncStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var total, active, del int
	var maxUpdated sql.NullInt64

	_ = d.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&total)
	_ = d.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE deleted_at IS NULL").Scan(&active)
	_ = d.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE deleted_at IS NOT NULL").Scan(&del)
	_ = d.db.QueryRow("SELECT MAX(updated_at) FROM transactions").Scan(&maxUpdated)

	return models.SyncStats{
		DeviceID:       d.deviceID,
		TotalRecords:   total,
		ActiveRecords:  active,
		DeletedRecords: del,
		LastUpdated:    maxUpdated.Int64,
	}
}

func (d *DB) ExportSyncPackage() (*models.SyncPackage, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query("SELECT " + txColumns + " FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records, err := scanTransactions(rows)
	if err != nil {
		return nil, err
	}

	return &models.SyncPackage{
		DeviceID:  d.deviceID,
		Timestamp: time.Now().UnixMilli(),
		Records:   records,
	}, nil
}

func (d *DB) ImportAndMergeSyncPackage(pkg *models.SyncPackage) (inserted, updated, skipped int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return 0, 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, remote := range pkg.Records {
		var localUpdated int64
		err = tx.QueryRow("SELECT updated_at FROM transactions WHERE id = ?", remote.ID).Scan(&localUpdated)
		if err == sql.ErrNoRows {
			// Record doesn't exist locally -> Insert
			insertQuery := `
			INSERT INTO transactions (
				id, date, type, benf, exp_cat, amount_paisa, account, notes,
				created_at, updated_at, deleted_at, device_id, version
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`
			_, err = tx.Exec(insertQuery,
				remote.ID, remote.Date, string(remote.Type), remote.Benf,
				remote.ExpCat, remote.Amount, string(remote.Account),
				remote.Notes, remote.CreatedAt, remote.UpdatedAt,
				remote.DeletedAt, remote.DeviceID, remote.Version,
			)
			if err != nil {
				return 0, 0, 0, err
			}
			inserted++
		} else if err == nil {
			// Record exists -> Compare timestamps (Last-Write-Wins)
			if remote.UpdatedAt > localUpdated {
				updateQuery := `
				UPDATE transactions SET
					date = ?,
					type = ?,
					benf = ?,
					exp_cat = ?,
					amount_paisa = ?,
					account = ?,
					notes = ?,
					updated_at = ?,
					deleted_at = ?,
					device_id = ?,
					version = ?
				WHERE id = ?
				`
				_, err = tx.Exec(updateQuery,
					remote.Date, string(remote.Type), remote.Benf,
					remote.ExpCat, remote.Amount, string(remote.Account),
					remote.Notes, remote.UpdatedAt, remote.DeletedAt,
					remote.DeviceID, remote.Version, remote.ID,
				)
				if err != nil {
					return 0, 0, 0, err
				}
				updated++
			} else {
				skipped++
			}
		} else {
			return 0, 0, 0, err
		}
	}

	err = tx.Commit()
	return inserted, updated, skipped, err
}

func (d *DB) ExportCSV() (string, error) {
	txs, err := d.GetTransactions(QueryFilter{})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("id,date,type,beneficiary,category,amount,account,notes\n")
	for _, t := range txs {
		benf := strings.ReplaceAll(t.Benf, `"`, `""`)
		cat := strings.ReplaceAll(t.ExpCat, `"`, `""`)
		notes := strings.ReplaceAll(t.Notes, `"`, `""`)
		sb.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%d.%02d,%s,\"%s\"\n",
			t.ID, t.Date, t.Type, benf, cat, t.Amount/100, t.Amount%100, t.Account, notes,
		))
	}
	return sb.String(), nil
}
