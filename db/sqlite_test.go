package db

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shakilnwz/wallet/core/models"
)

func TestSQLiteDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_wallet.db")

	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer database.Close()

	if database.DeviceID() == "" {
		t.Errorf("DeviceID is empty")
	}

	// 1. Add Income
	inc, err := database.AddTransaction(models.Transaction{
		Date:    "2026-10-01",
		Type:    models.TypeIncome,
		Benf:    "Salary",
		ExpCat:  "Salary",
		Amount:  50000.0,
		Account: models.AccountBank,
		Notes:   "October salary",
	})
	if err != nil {
		t.Fatalf("AddTransaction income failed: %v", err)
	}
	if inc.ID == "" {
		t.Errorf("Expected generated UUID")
	}

	// 2. Add Expense 1 (Bank)
	_, err = database.AddTransaction(models.Transaction{
		Date:    "2026-10-02",
		Type:    models.TypeExpense,
		Benf:    "Rent",
		ExpCat:  "Housing",
		Amount:  15000.0,
		Account: models.AccountBank,
	})
	if err != nil {
		t.Fatalf("AddTransaction rent failed: %v", err)
	}

	// 3. Add Expense 2 (Cash)
	food, err := database.AddTransaction(models.Transaction{
		Date:    "2026-10-03",
		Type:    models.TypeExpense,
		Benf:    "Grocery",
		ExpCat:  "Food",
		Amount:  5000.0,
		Account: models.AccountCash,
	})
	if err != nil {
		t.Fatalf("AddTransaction food failed: %v", err)
	}

	// 4. Test Account Summary
	summary, err := database.GetAccountSummary("2026-10")
	if err != nil {
		t.Fatalf("GetAccountSummary failed: %v", err)
	}

	if summary.TotalIncome != 50000.0 {
		t.Errorf("Expected TotalIncome 50000.0, got %f", summary.TotalIncome)
	}
	if summary.TotalExpense != 20000.0 {
		t.Errorf("Expected TotalExpense 20000.0, got %f", summary.TotalExpense)
	}
	if summary.TotalBalance != 30000.0 {
		t.Errorf("Expected TotalBalance 30000.0, got %f", summary.TotalBalance)
	}
	if summary.BankBalance != 35000.0 {
		t.Errorf("Expected BankBalance 35000.0, got %f", summary.BankBalance)
	}
	if summary.CashBalance != -5000.0 {
		t.Errorf("Expected CashBalance -5000.0, got %f", summary.CashBalance)
	}
	if !summary.IsBalanced {
		t.Errorf("Expected IsBalanced true, got false (diff=%f)", summary.NeedBalancing)
	}

	// 5. Soft Delete
	err = database.SoftDeleteTransaction(food.ID)
	if err != nil {
		t.Fatalf("SoftDelete failed: %v", err)
	}

	afterDelete, err := database.GetAccountSummary("2026-10")
	if err != nil {
		t.Fatalf("GetAccountSummary after delete failed: %v", err)
	}
	if afterDelete.TotalExpense != 15000.0 {
		t.Errorf("Expected TotalExpense 15000.0 after soft delete, got %f", afterDelete.TotalExpense)
	}

	// 6. Test Sync LWW
	syncPkg, err := database.ExportSyncPackage()
	if err != nil {
		t.Fatalf("ExportSyncPackage failed: %v", err)
	}
	if len(syncPkg.Records) != 3 {
		t.Errorf("Expected 3 records in sync package (including tombstone), got %d", len(syncPkg.Records))
	}

	// Simulate remote peer merging
	peerDbPath := filepath.Join(tempDir, "peer_wallet.db")
	peerDb, err := Open(peerDbPath)
	if err != nil {
		t.Fatalf("Open peer failed: %v", err)
	}
	defer peerDb.Close()

	inserted, updated, skipped, err := peerDb.ImportAndMergeSyncPackage(syncPkg)
	if err != nil {
		t.Fatalf("ImportAndMergeSyncPackage failed: %v", err)
	}
	if inserted != 3 || updated != 0 || skipped != 0 {
		t.Errorf("Expected 3 inserted, got %d ins, %d upd, %d skip", inserted, updated, skipped)
	}

	// Repeat import -> should all be skipped
	_, _, skipped2, err := peerDb.ImportAndMergeSyncPackage(syncPkg)
	if err != nil {
		t.Fatalf("Second ImportAndMergeSyncPackage failed: %v", err)
	}
	if skipped2 != 3 {
		t.Errorf("Expected 3 skipped on repeat import, got %d", skipped2)
	}

	// Test LWW update from peer
	now := time.Now().UnixMilli() + 1000
	peerPkg := &models.SyncPackage{
		DeviceID:  "peer-device",
		Timestamp: now,
		Records: []models.Transaction{
			{
				ID:        inc.ID,
				Date:      "2026-10-01",
				Type:      models.TypeIncome,
				Benf:      "Salary (Bonus Included)",
				ExpCat:    "Salary",
				Amount:    55000.0,
				Account:   models.AccountBank,
				Notes:     "Updated October salary",
				CreatedAt: inc.CreatedAt,
				UpdatedAt: now,
				DeletedAt: nil,
				DeviceID:  "peer-device",
				Version:   2,
			},
		},
	}

	ins3, upd3, _, err := database.ImportAndMergeSyncPackage(peerPkg)
	if err != nil {
		t.Fatalf("Import remote updated failed: %v", err)
	}
	if ins3 != 0 || upd3 != 1 {
		t.Errorf("Expected 1 updated, got %d ins, %d upd", ins3, upd3)
	}

	txs, err := database.GetTransactions(QueryFilter{Month: "2026-10"})
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}
	var foundUpdated bool
	for _, tx := range txs {
		if tx.ID == inc.ID && tx.Amount == 55000.0 && tx.Benf == "Salary (Bonus Included)" {
			foundUpdated = true
			break
		}
	}
	if !foundUpdated {
		t.Errorf("Remote LWW update was not applied properly")
	}
}
