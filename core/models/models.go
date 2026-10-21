package models

type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

type AccountType string

const (
	AccountBank AccountType = "Account"
	AccountCash AccountType = "Cash"
)

type Transaction struct {
	ID        string          `json:"id"`
	Date      string          `json:"date"`       // YYYY-MM-DD
	Type      TransactionType `json:"type"`       // income | expense
	Benf      string          `json:"benf"`       // Beneficiary / Payee / Source
	ExpCat    string          `json:"exp_cat"`    // Category
	Amount    float64         `json:"amount"`     // BDT (৳)
	Account   AccountType     `json:"account"`    // Account | Cash
	Notes     string          `json:"notes"`      // Optional remarks
	CreatedAt int64           `json:"created_at"` // Epoch ms
	UpdatedAt int64           `json:"updated_at"` // Epoch ms
	DeletedAt *int64          `json:"deleted_at"` // Epoch ms for tombstone
	DeviceID  string          `json:"device_id"`  // Originating device UUID
	Version   int             `json:"version"`    // Incremented on local edit
}

type AccountSummary struct {
	TotalIncome   float64 `json:"total_income"`
	TotalExpense  float64 `json:"total_expense"`
	TotalBalance  float64 `json:"total_balance"`
	BankBalance   float64 `json:"bank_balance"`
	CashBalance   float64 `json:"cash_balance"`
	NeedBalancing float64 `json:"need_balancing"`
	IsBalanced    bool    `json:"is_balanced"`
}

type CategoryBreakdown struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type MonthlySummary struct {
	Month      string              `json:"month"` // YYYY-MM
	Income     float64             `json:"income"`
	Expense    float64             `json:"expense"`
	Balance    float64             `json:"balance"`
	Categories []CategoryBreakdown `json:"categories"`
}

type SyncStats struct {
	DeviceID       string `json:"device_id"`
	TotalRecords   int    `json:"total_records"`
	ActiveRecords  int    `json:"active_records"`
	DeletedRecords int    `json:"deleted_records"`
	LastUpdated    int64  `json:"last_updated"`
}

type SyncPackage struct {
	DeviceID  string        `json:"device_id"`
	Timestamp int64         `json:"timestamp"`
	Records   []Transaction `json:"records"`
}
