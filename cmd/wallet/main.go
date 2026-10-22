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

		walletUI.Invalidate = w.Invalidate
		loopErr := runEventLoop(w, walletUI)

		// app.Main never returns (on Android too), so deferred calls in main
		// would never run: shut down explicitly before exiting.
		walletUI.Shutdown()
		if err := database.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
		if loopErr != nil {
			log.Fatalf("Window event loop error: %v", loopErr)
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
