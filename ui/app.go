package ui

import (
	"fmt"
	"image"
	"slices"
	stdsync "sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/db"
	"github.com/shakilnwz/wallet/sync"
	"github.com/shakilnwz/wallet/ui/components"
	"github.com/shakilnwz/wallet/ui/theme"
)

const syncPort = 8080

type App struct {
	Theme    *theme.Theme
	Database *db.DB

	CurrentTab    components.Tab
	SelectedMonth string // YYYY-MM or 'all'
	Months        []string

	Transactions   []models.Transaction
	AccountSummary *models.AccountSummary
	MonthlySummary *models.MonthlySummary
	SyncStats      models.SyncStats
	LocalIP        string

	// Components
	Header        components.Header
	NavDrawer     components.NavDrawer
	MonthSelector components.MonthSelector
	HeroCard      components.HeroCard
	TxList        components.TxList
	BottomNav     components.BottomNav
	QuickAddSheet components.QuickAddSheet
	AnalyticsView components.AnalyticsView
	SyncView      components.SyncView
	SettingsView  components.SettingsView
	Snackbar      components.Snackbar

	// Home screen shortcuts
	HomeViewAllBtn  widget.Clickable
	HomeQuickAddBtn widget.Clickable

	// Sync
	SyncServer *sync.Server
	SyncClient *sync.Client

	// Invalidate wakes the event loop for a redraw. It is safe to call from
	// any goroutine. Set by main before the first frame.
	Invalidate func()

	mu      stdsync.Mutex
	pending []func() // UI updates queued by background goroutines
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
	app.QuickAddSheet.Categories = app.categoriesFor
	app.SyncServer.OnMerge = func() {
		app.post(func() {
			app.ReloadData()
			app.Snackbar.Show("Received changes from another device", false, 3*time.Second, "", nil)
		})
	}

	app.ReloadData()
	return app
}

// post queues f to run on the UI goroutine at the start of the next frame.
func (a *App) post(f func()) {
	a.mu.Lock()
	a.pending = append(a.pending, f)
	a.mu.Unlock()
	if a.Invalidate != nil {
		a.Invalidate()
	}
}

func (a *App) runPending() {
	a.mu.Lock()
	fs := a.pending
	a.pending = nil
	a.mu.Unlock()
	for _, f := range fs {
		f()
	}
}

func (a *App) showError(prefix string, err error) {
	a.Snackbar.Show(fmt.Sprintf("%s: %v", prefix, err), true, 6*time.Second, "", nil)
}

// categoriesFor lists the user's most used categories for a type first,
// followed by any built-in defaults not already present.
func (a *App) categoriesFor(t models.TransactionType) []string {
	recent, err := a.Database.GetRecentCategories(t)
	if err != nil {
		a.showError("Couldn't load categories", err)
	}
	cats := slices.Clone(recent)
	for _, c := range components.DefaultCategories(t) {
		if !slices.Contains(cats, c) {
			cats = append(cats, c)
		}
	}
	return cats
}

// ReloadData refreshes everything shown on screen from the database. On
// error the previous values are kept and the error is shown.
func (a *App) ReloadData() {
	var firstErr error
	keep := func(err error) bool {
		if err != nil && firstErr == nil {
			firstErr = err
		}
		return err == nil
	}

	if months, err := a.Database.GetAvailableMonths(); keep(err) {
		a.Months = months
	}

	monthParam := a.SelectedMonth
	if a.SelectedMonth == "all" {
		monthParam = ""
	}

	if txs, err := a.Database.GetTransactions(db.QueryFilter{Month: monthParam}); keep(err) {
		a.Transactions = txs
	}
	if s, err := a.Database.GetAccountSummary(monthParam); keep(err) {
		a.AccountSummary = s
	}
	if s, err := a.Database.GetMonthlySummary(a.SelectedMonth); keep(err) {
		a.MonthlySummary = s
	}
	a.SyncStats = a.Database.GetSyncStats()
	a.LocalIP = sync.LocalIP()

	if firstErr != nil {
		a.showError("Couldn't load data", firstErr)
	}
	if a.Invalidate != nil {
		a.Invalidate()
	}
}

func (a *App) saveTransaction(tx models.Transaction) {
	var err error
	if tx.ID != "" {
		err = a.Database.UpdateTransaction(tx)
	} else {
		_, err = a.Database.AddTransaction(tx)
	}
	if err != nil {
		a.showError("Couldn't save", err)
		return
	}
	// Jump to the saved entry's month so it's visible.
	if a.SelectedMonth != "all" && len(tx.Date) >= 7 {
		a.SelectedMonth = tx.Date[:7]
	}
	a.ReloadData()
}

func (a *App) deleteTransaction(tx models.Transaction) {
	if err := a.Database.SoftDeleteTransaction(tx.ID); err != nil {
		a.showError("Couldn't delete", err)
		return
	}
	a.ReloadData()
	a.Snackbar.Show("Deleted "+tx.Benf, false, 5*time.Second, "Undo", func() {
		if err := a.Database.RestoreTransaction(tx.ID); err != nil {
			a.showError("Couldn't restore", err)
		}
		a.ReloadData()
	})
}

// startSync runs a two-way sync in the background so the UI stays responsive.
func (a *App) startSync(peer string) {
	if a.SyncView.Syncing {
		return
	}
	a.SyncView.Syncing = true
	a.SyncView.SetStatus("Syncing with "+peer+"...", false)
	if a.Invalidate != nil {
		a.Invalidate()
	}

	go func() {
		pIns, pUpd, sIns, sUpd, err := a.SyncClient.FullTwoWaySync(peer, a.Database)
		a.post(func() {
			a.SyncView.Syncing = false
			if err != nil {
				a.SyncView.SetStatus(fmt.Sprintf("Sync failed: %v", err), true)
				return
			}
			a.SyncView.SetStatus(fmt.Sprintf("Synced. Received %d new, %d updated. Sent %d new, %d updated.", pIns, pUpd, sIns, sUpd), false)
			a.ReloadData()
		})
	}()
}

func (a *App) toggleServer() {
	if a.SyncServer.IsRunning() {
		if err := a.SyncServer.Stop(); err != nil {
			a.showError("Couldn't stop server", err)
		}
	} else {
		if _, err := a.SyncServer.Start(syncPort); err != nil {
			a.showError("Couldn't start server", err)
		}
	}
	a.LocalIP = sync.LocalIP()
	if a.Invalidate != nil {
		a.Invalidate()
	}
}

func (a *App) deleteAll() {
	n, err := a.Database.ClearAllData()
	if err != nil {
		a.SyncView.SetStatus(fmt.Sprintf("Delete failed: %v", err), true)
	} else {
		a.SyncView.SetStatus(fmt.Sprintf("Deleted %d transactions. Synced devices will delete them at their next sync.", n), false)
	}
	a.ReloadData()
}

func (a *App) selectTab(t components.Tab) {
	a.CurrentTab = t
	if t == components.TabSync || t == components.TabSettings {
		a.SyncStats = a.Database.GetSyncStats()
		a.LocalIP = sync.LocalIP()
	}
	if a.Invalidate != nil {
		a.Invalidate()
	}
}

// Shutdown stops background services. Call it before closing the database.
func (a *App) Shutdown() {
	_ = a.SyncServer.Stop()
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	a.runPending()

	// 1. Fill Background with OLED Black
	paint.Fill(gtx.Ops, theme.ColorBackground)
	gtx.Constraints.Min = gtx.Constraints.Max

	// 2. Main screen, with the drawer, snackbar, and quick-add sheet layered on top.
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(a.layoutMain),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Navigation Drawer (slides out smoothly from left)
			return a.NavDrawer.Layout(gtx, a.Theme, a.CurrentTab, a.LocalIP, a.SyncStats.ActiveRecords, a.selectTab)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Quick Add Sheet (slides up smoothly from bottom)
			return a.QuickAddSheet.Layout(gtx, a.Theme, a.saveTransaction)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// Sit just above the bottom navigation bar.
			return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(90)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return a.Snackbar.Layout(gtx, a.Theme)
				})
			})
		}),
	)
}

func (a *App) layoutMain(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Top App Bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			title := "Wallet"
			sub := "Dashboard"
			actionLabel := ""
			switch a.CurrentTab {
			case components.TabLedger:
				title = "Ledger"
				sub = "All Transactions"
			case components.TabInsights:
				title = "Insights"
				sub = "Spending Analytics"
			case components.TabSync:
				title = "P2P Sync"
				sub = "Local Wi-Fi Network"
			case components.TabSettings:
				title = "Data & Settings"
				sub = "Database Management"
			}

			return a.Header.Layout(gtx, a.Theme,
				title,
				sub,
				actionLabel,
				a.NavDrawer.Toggle,
				a.ReloadData,
				nil,
			)
		}),

		// Main Active Screen Content (Full Width)
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			switch a.CurrentTab {
			case components.TabHome:
				return a.layoutHomeScreen(gtx)

			case components.TabLedger:
				return a.layoutLedgerScreen(gtx)

			case components.TabInsights:
				return a.layoutInsightsScreen(gtx)

			case components.TabSync:
				return a.SyncView.Layout(gtx, a.Theme, a.LocalIP, a.SyncServer.IsRunning(), a.SyncServer.Port(),
					a.toggleServer,
					a.startSync,
				)

			case components.TabSettings:
				return a.SettingsView.Layout(gtx, a.Theme, a.SyncStats,
					a.deleteAll,
				)
			}
			return layout.Dimensions{}
		}),

		// Bottom Thumb Navigation (Full Width)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.BottomNav.Layout(gtx, a.Theme, a.CurrentTab, a.selectTab, a.QuickAddSheet.Open, a.NavDrawer.Toggle)
		}),
	)
}

// layoutHomeScreen renders the clean, decluttered dashboard
func (a *App) layoutHomeScreen(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	// Handle Home shortcut clicks
	if a.HomeViewAllBtn.Clicked(gtx) {
		a.selectTab(components.TabLedger)
		gtx.Execute(op.InvalidateCmd{})
	}
	if a.HomeQuickAddBtn.Clicked(gtx) {
		a.QuickAddSheet.Open()
		gtx.Execute(op.InvalidateCmd{})
	}

	// Show only the 5 most recent transactions on the home screen to declutter
	recentTxs := a.Transactions
	if len(recentTxs) > 5 {
		recentTxs = recentTxs[:5]
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Hero Balance Card
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.HeroCard.Layout(gtx, a.Theme, a.AccountSummary, theme.FormatMonthDisplay(a.SelectedMonth))
		}),

		// Quick Action Shortcuts Bar (Touch-Friendly)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(a.Theme.Theme, &a.HomeQuickAddBtn, "+ Add Entry")
						btn.Background = theme.ColorPrimary
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(12)
						btn.TextSize = unit.Sp(13)
						btn.Inset = layout.Inset{Top: unit.Dp(11), Bottom: unit.Dp(11)}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(a.Theme.Theme, &a.HomeViewAllBtn, "Full History →")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(12)
						btn.TextSize = unit.Sp(13)
						btn.Inset = layout.Inset{Top: unit.Dp(11), Bottom: unit.Dp(11)}
						return btn.Layout(gtx)
					}),
				)
			})
		}),

		// Recent Activity Section Header
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Left: unit.Dp(18), Right: unit.Dp(18), Top: unit.Dp(4), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(a.Theme.Theme, "RECENT ACTIVITY")
						lbl.Color = theme.ColorTextMuted
						lbl.TextSize = unit.Sp(11)
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(a.Theme.Theme, &a.HomeViewAllBtn, "See All")
						btn.Background = theme.ColorSurface
						btn.Color = theme.ColorPrimary
						btn.CornerRadius = unit.Dp(6)
						btn.TextSize = unit.Sp(11)
						btn.Inset = layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(8), Right: unit.Dp(8)}
						return btn.Layout(gtx)
					}),
				)
			})
		}),

		// Recent Transactions List
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.TxList.Layout(gtx, a.Theme, recentTxs, a.QuickAddSheet.OpenEdit, a.deleteTransaction)
		}),
	)
}

// layoutLedgerScreen renders the complete transaction register with month filter
func (a *App) layoutLedgerScreen(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.MonthSelector.Layout(gtx, a.Theme, a.SelectedMonth, a.Months, func(m string) {
				a.SelectedMonth = m
				a.ReloadData()
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.TxList.Layout(gtx, a.Theme, a.Transactions, a.QuickAddSheet.OpenEdit, a.deleteTransaction)
		}),
	)
}

// layoutInsightsScreen renders spending analysis & account distribution
func (a *App) layoutInsightsScreen(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.MonthSelector.Layout(gtx, a.Theme, a.SelectedMonth, a.Months, func(m string) {
				a.SelectedMonth = m
				a.ReloadData()
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return a.AnalyticsView.Layout(gtx, a.Theme, a.MonthlySummary, a.AccountSummary)
		}),
	)
}
