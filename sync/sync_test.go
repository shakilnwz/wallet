package sync

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/db"
)

func TestSyncWithOutdatedPeer(t *testing.T) {
	// An old build only serves /sync/*, so every /sync/v2 path 404s.
	old := httptest.NewServer(http.NotFoundHandler())
	defer old.Close()

	database, err := db.Open(filepath.Join(t.TempDir(), "d.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	_, _, _, _, err = NewClient().FullTwoWaySync(old.URL, database)
	if !errors.Is(err, errPeerOutdated) {
		t.Fatalf("got %v, want errPeerOutdated", err)
	}
}

func TestSyncServerClient(t *testing.T) {
	tempDir := t.TempDir()

	// Device 1 (Desktop)
	db1, err := db.Open(filepath.Join(tempDir, "device1.db"))
	if err != nil {
		t.Fatalf("Open db1 failed: %v", err)
	}
	defer db1.Close()

	_, err = db1.AddTransaction(models.Transaction{
		Date:    "2026-10-01",
		Type:    models.TypeIncome,
		Benf:    "Freelance Project",
		ExpCat:  "Freelance",
		Amount:  40000_00,
		Account: models.AccountBank,
	})
	if err != nil {
		t.Fatalf("AddTransaction on db1 failed: %v", err)
	}

	// Device 2 (Phone)
	db2, err := db.Open(filepath.Join(tempDir, "device2.db"))
	if err != nil {
		t.Fatalf("Open db2 failed: %v", err)
	}
	defer db2.Close()

	_, err = db2.AddTransaction(models.Transaction{
		Date:    "2026-10-02",
		Type:    models.TypeExpense,
		Benf:    "Coffee",
		ExpCat:  "Food & Dining",
		Amount:  250_00,
		Account: models.AccountCash,
	})
	if err != nil {
		t.Fatalf("AddTransaction on db2 failed: %v", err)
	}

	// Start sync server on Device 1
	server1 := NewServer(db1)
	var merged atomic.Bool
	server1.OnMerge = func() { merged.Store(true) }
	port1, err := server1.Start(0) // Random port
	if err != nil {
		t.Fatalf("Start server1 failed: %v", err)
	}
	defer server1.Stop()

	// Client on Device 2 syncs with Device 1
	client2 := NewClient()
	peerAddr := fmt.Sprintf("127.0.0.1:%d", port1)

	peerID, err := client2.Ping(peerAddr)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	if peerID != db1.DeviceID() {
		t.Errorf("Expected peer ID %s, got %s", db1.DeviceID(), peerID)
	}

	// Run two-way sync
	pulledIns, pulledUpd, pushedIns, pushedUpd, err := client2.FullTwoWaySync(peerAddr, db2)
	if err != nil {
		t.Fatalf("FullTwoWaySync failed: %v", err)
	}

	if pulledIns != 1 || pulledUpd != 0 {
		t.Errorf("Expected 1 pulled insert on db2, got %d ins, %d upd", pulledIns, pulledUpd)
	}
	if pushedIns != 1 || pushedUpd != 0 {
		t.Errorf("Expected 1 pushed insert on db1, got %d ins, %d upd", pushedIns, pushedUpd)
	}

	// Check both databases now have both records (2 records each)
	txs1, _ := db1.GetTransactions(db.QueryFilter{})
	txs2, _ := db2.GetTransactions(db.QueryFilter{})

	if len(txs1) != 2 {
		t.Errorf("Expected 2 transactions in db1, got %d", len(txs1))
	}
	if len(txs2) != 2 {
		t.Errorf("Expected 2 transactions in db2, got %d", len(txs2))
	}
	if !merged.Load() {
		t.Errorf("Expected OnMerge to fire after the peer pushed records")
	}
}
