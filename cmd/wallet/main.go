package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/shakilnwz/wallet/db"
	"github.com/shakilnwz/wallet/ui"
)

func main() {
	// Initialize local SQLite database
	database, err := db.Open("")
	if err != nil {
		log.Fatalf("Failed to initialize SQLite database: %v", err)
	}
	defer database.Close()

	// Initialize UI application
	walletUI := ui.NewApp(database)

	go func() {
		// Standard mobile phone viewport dimensions (390 x 844 dp)
		w := new(app.Window)
		w.Option(
			app.Title("Wallet"),
			app.Size(unit.Dp(390), unit.Dp(844)),
			app.MinSize(unit.Dp(320), unit.Dp(568)),
		)

		if err := runEventLoop(w, walletUI); err != nil {
			log.Fatalf("Window event loop error: %v", err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func runEventLoop(w *app.Window, walletUI *ui.App) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			walletUI.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
