package db

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
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
		// Determine suitable directory
		baseDir, err := os.UserConfigDir()
		if err != nil || baseDir == "" {
			baseDir = "."
		}
		appDir := filepath.Join(baseDir, "wallet")
		_ = os.MkdirAll(appDir, 0755)
		dbPath = filepath.Join(appDir, "wallet.db")
	}

	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Enable WAL mode & foreign keys for performance and crash resistance
	_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL;")
	_, _ = sqlDB.Exec("PRAGMA busy_timeout=5000;")

	d := &DB{db: sqlDB}
	if err := d.initSchema(); err != nil {
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

func (d *DB) initSchema() error {
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

	CREATE TABLE IF NOT EXISTS meta (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`
	_, err := d.db.Exec(schema)
	return err
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
	sb.WriteString("SELECT id, date, type, benf, exp_cat, amount, account, COALESCE(notes, ''), created_at, updated_at, deleted_at, device_id, version ")
	sb.WriteString("FROM transactions WHERE deleted_at IS NULL ")

	var args []any

	if filter.Month != "" && filter.Month != "all" {
		sb.WriteString("AND date LIKE ? ")
		args = append(args, filter.Month+"%")
	}

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

	return results, nil
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
		id, date, type, benf, exp_cat, amount, account, notes,
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
		amount = ?,
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

func (d *DB) ClearAllData() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, err := d.db.Exec("DELETE FROM transactions")
	return err
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
		COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as inc,
		COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as exp,
		COALESCE(SUM(CASE WHEN account = 'Account' THEN (CASE WHEN type = 'income' THEN amount ELSE -amount END) ELSE 0 END), 0) as bank,
		COALESCE(SUM(CASE WHEN account = 'Cash' THEN (CASE WHEN type = 'income' THEN amount ELSE -amount END) ELSE 0 END), 0) as cash
	FROM transactions
	WHERE deleted_at IS NULL
	`)

	var args []any
	if month != "" && month != "all" {
		sb.WriteString(" AND date LIKE ?")
		args = append(args, month+"%")
	}

	var inc, exp, bank, cash float64
	err := d.db.QueryRow(sb.String(), args...).Scan(&inc, &exp, &bank, &cash)
	if err != nil {
		return nil, err
	}

	totalBal := inc - exp
	diff := math.Round((totalBal-(bank+cash))*100) / 100

	return &models.AccountSummary{
		TotalIncome:   inc,
		TotalExpense:  exp,
		TotalBalance:  totalBal,
		BankBalance:   bank,
		CashBalance:   cash,
		NeedBalancing: diff,
		IsBalanced:    math.Abs(diff) < 0.01,
	}, nil
}

func (d *DB) GetMonthlySummary(month string) (*models.MonthlySummary, error) {
	accSummary, err := d.GetAccountSummary(month)
	if err != nil {
		return nil, err
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	query := `
	SELECT exp_cat, SUM(amount) as total
	FROM transactions
	WHERE deleted_at IS NULL AND type = 'expense' AND date LIKE ?
	GROUP BY exp_cat
	ORDER BY total DESC
	`

	rows, err := d.db.Query(query, month+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryBreakdown
	for rows.Next() {
		var cat string
		var amount float64
		if err := rows.Scan(&cat, &amount); err != nil {
			return nil, err
		}
		pct := 0.0
		if accSummary.TotalExpense > 0 {
			pct = math.Round((amount/accSummary.TotalExpense)*1000) / 10
		}
		categories = append(categories, models.CategoryBreakdown{
			Category:   cat,
			Amount:     amount,
			Percentage: pct,
		})
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

	var months []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err == nil {
			months = append(months, m)
		}
	}
	return months, nil
}

func (d *DB) GetRecentCategories() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT exp_cat, COUNT(*) as cnt
		FROM transactions
		WHERE deleted_at IS NULL
		GROUP BY exp_cat
		ORDER BY cnt DESC
		LIMIT 15
	`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var c string
		var cnt int
		if err := rows.Scan(&c, &cnt); err == nil {
			cats = append(cats, c)
		}
	}
	return cats
}

func (d *DB) GetRecentBeneficiaries() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.Query(`
		SELECT benf, COUNT(*) as cnt
		FROM transactions
		WHERE deleted_at IS NULL
		GROUP BY benf
		ORDER BY cnt DESC
		LIMIT 15
	`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var benfs []string
	for rows.Next() {
		var b string
		var cnt int
		if err := rows.Scan(&b, &cnt); err == nil {
			benfs = append(benfs, b)
		}
	}
	return benfs
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

	rows, err := d.db.Query("SELECT id, date, type, benf, exp_cat, amount, account, COALESCE(notes, ''), created_at, updated_at, deleted_at, device_id, version FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var delAt sql.NullInt64
		if err := rows.Scan(&tx.ID, &tx.Date, &tx.Type, &tx.Benf, &tx.ExpCat, &tx.Amount, &tx.Account, &tx.Notes, &tx.CreatedAt, &tx.UpdatedAt, &delAt, &tx.DeviceID, &tx.Version); err == nil {
			if delAt.Valid {
				tx.DeletedAt = &delAt.Int64
			}
			records = append(records, tx)
		}
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
				id, date, type, benf, exp_cat, amount, account, notes,
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
					amount = ?,
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
		sb.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%.2f,%s,\"%s\"\n",
			t.ID, t.Date, t.Type, benf, cat, t.Amount, t.Account, notes,
		))
	}
	return sb.String(), nil
}
