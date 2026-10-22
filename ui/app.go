package ui

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/db"
	"github.com/shakilnwz/wallet/sync"
	"github.com/shakilnwz/wallet/ui/components"
	"github.com/shakilnwz/wallet/ui/theme"
)

type App struct {
	Theme    *theme.Theme
	Database *db.DB

	CurrentTab    components.Tab
	SelectedMonth string // YYYY-MM or 'all'
	Months        []string

	Transactions   []models.Transaction
	AccountSummary *models.AccountSummary
	MonthlySummary *models.MonthlySummary

	// Components
	Header        components.Header
	MonthSelector components.MonthSelector
	HeroCard      components.HeroCard
	TxList        components.TxList
	BottomNav     components.BottomNav
	QuickAddSheet components.QuickAddSheet
	AnalyticsView components.AnalyticsView
	SyncView      components.SyncView

	// Sync
	SyncServer *sync.Server
	SyncClient *sync.Client
}

func NewApp(database *db.DB) *App {
	th := theme.NewTheme()
	curMonth := time.Now().Format("2006-01")

	app := &App{
		Theme:         th,
		Database:      database,
		CurrentTab:    components.TabHome,
		SelectedMonth: curMonth,
		SyncServer:    sync.NewServer(database),
		SyncClient:    sync.NewClient(),
	}

	app.ReloadData()
	return app
}

func (a *App) ReloadData() {
	months, _ := a.Database.GetAvailableMonths()
	a.Months = months

	if a.SelectedMonth != "all" && len(months) > 0 {
		var found bool
		for _, m := range months {
			if m == a.SelectedMonth {
				found = true
				break
			}
		}
		if !found && a.SelectedMonth == "" {
			a.SelectedMonth = months[0]
		}
	}

	monthParam := a.SelectedMonth
	if a.SelectedMonth == "all" {
		monthParam = ""
	}

	txs, _ := a.Database.GetTransactions(db.QueryFilter{
		Month: monthParam,
	})
	a.Transactions = txs

	accSummary, _ := a.Database.GetAccountSummary(monthParam)
	a.AccountSummary = accSummary

	monSummary, _ := a.Database.GetMonthlySummary(a.SelectedMonth)
	a.MonthlySummary = monSummary
}

// Shutdown stops background services. Call it before closing the database.
func (a *App) Shutdown() {
	_ = a.SyncServer.Stop()
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	// 1. Fill Background with OLED Black
	paint.Fill(gtx.Ops, theme.ColorBackground)

	// 2. Main Flex Layout (Header, Content, Bottom Navigation)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.Header.Layout(gtx, a.Theme,
				func() { // On Refresh
					a.ReloadData()
				},
				func() { // On Open Sync
					a.CurrentTab = components.TabSync
				},
			)
		}),

		// Month Selector Carousel (Shown in Home, Ledger, and Insights tabs)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if a.CurrentTab == components.TabSync {
				return layout.Dimensions{}
			}
			return a.MonthSelector.Layout(gtx, a.Theme, a.SelectedMonth, a.Months, func(m string) {
				a.SelectedMonth = m
				a.ReloadData()
			})
		}),

		// Main Active Tab Content
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			switch a.CurrentTab {
			case components.TabHome:
				// Home: Hero Balance Card + Recent Transactions
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.HeroCard.Layout(gtx, a.Theme, a.AccountSummary, theme.FormatMonthDisplay(a.SelectedMonth))
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.TxList.Layout(gtx, a.Theme, a.Transactions, func(id string) {
							_ = a.Database.SoftDeleteTransaction(id)
							a.ReloadData()
						})
					}),
				)

			case components.TabLedger:
				// Ledger: Full History
				return a.TxList.Layout(gtx, a.Theme, a.Transactions, func(id string) {
					_ = a.Database.SoftDeleteTransaction(id)
					a.ReloadData()
				})

			case components.TabInsights:
				// Insights: Spending Bars & Distribution
				return a.AnalyticsView.Layout(gtx, a.Theme, a.MonthlySummary, a.AccountSummary)

			case components.TabSync:
				// Sync & DB Management
				stats := a.Database.GetSyncStats()
				return a.SyncView.Layout(gtx, a.Theme, stats, a.SyncServer.IsRunning(), a.SyncServer.Port(),
					func() { // Toggle Server
						if a.SyncServer.IsRunning() {
							_ = a.SyncServer.Stop()
						} else {
							_, _ = a.SyncServer.Start(8080)
						}
					},
					func(peerAddr string) { // Sync With Peer
						pIns, pUpd, sIns, sUpd, err := a.SyncClient.FullTwoWaySync(peerAddr, a.Database)
						if err != nil {
							a.SyncView.StatusMsg = fmt.Sprintf("Sync failed: %v", err)
						} else {
							a.SyncView.StatusMsg = fmt.Sprintf("Synced! Pulled %d+%d, Pushed %d+%d", pIns, pUpd, sIns, sUpd)
							a.ReloadData()
						}
					},
					func() { // Delete all transactions (propagates to peers on sync)
						n, err := a.Database.ClearAllData()
						if err != nil {
							a.SyncView.StatusMsg = fmt.Sprintf("Delete failed: %v", err)
						} else {
							a.SyncView.StatusMsg = fmt.Sprintf("Deleted %d transactions. Peers will delete them at next sync.", n)
						}
						a.ReloadData()
					},
				)
			}
			return layout.Dimensions{}
		}),

		// Bottom Thumb Navigation
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.BottomNav.Layout(gtx, a.Theme, a.CurrentTab,
				func(t components.Tab) { // Tab change
					a.CurrentTab = t
				},
				func() { // Open Quick Add
					a.QuickAddSheet.Open()
				},
			)
		}),

		// Modal Overlay (Quick Add Sheet)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.QuickAddSheet.Layout(gtx, a.Theme, func(tx models.Transaction) {
				_, _ = a.Database.AddTransaction(tx)
				a.ReloadData()
			})
		}),
	)
}
