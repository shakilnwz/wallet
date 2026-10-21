package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type Tab int

const (
	TabHome Tab = iota
	TabLedger
	TabInsights
	TabSync
)

type BottomNav struct {
	HomeBtn     widget.Clickable
	LedgerBtn   widget.Clickable
	InsightsBtn widget.Clickable
	SyncBtn     widget.Clickable
	AddBtn      widget.Clickable
}

func (b *BottomNav) Layout(
	gtx layout.Context,
	th *theme.Theme,
	currentTab Tab,
	onSelectTab func(Tab),
	onOpenAdd func(),
) layout.Dimensions {
	if b.HomeBtn.Clicked(gtx) {
		onSelectTab(TabHome)
	}
	if b.LedgerBtn.Clicked(gtx) {
		onSelectTab(TabLedger)
	}
	if b.InsightsBtn.Clicked(gtx) {
		onSelectTab(TabInsights)
	}
	if b.SyncBtn.Clicked(gtx) {
		onSelectTab(TabSync)
	}
	if b.AddBtn.Clicked(gtx) {
		onOpenAdd()
	}

	return layout.Inset{
		Bottom: unit.Dp(8),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return theme.Card(gtx, theme.ColorSurface, unit.Dp(24), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					// Home tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.HomeBtn, "Home", currentTab == TabHome)
					}),

					// Ledger tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.LedgerBtn, "Ledger", currentTab == TabLedger)
					}),

					// Center Prominent Add (+) Button
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &b.AddBtn, "+")
						btn.Background = theme.ColorPrimary
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(22)
						btn.TextSize = unit.Sp(22)
						btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(18), Right: unit.Dp(18)}
						return btn.Layout(gtx)
					}),

					// Insights tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.InsightsBtn, "Insights", currentTab == TabInsights)
					}),

					// Sync tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.SyncBtn, "Sync", currentTab == TabSync)
					}),
				)
			})
		})
	})
}

func (b *BottomNav) layoutTabBtn(gtx layout.Context, th *theme.Theme, click *widget.Clickable, title string, isActive bool) layout.Dimensions {
	col := theme.ColorTextMuted
	bg := theme.ColorSurface
	if isActive {
		col = theme.ColorPrimary
		bg = theme.ColorSurfaceElevated
	}

	btn := material.Button(th.Theme, click, title)
	btn.Background = bg
	btn.Color = col
	btn.CornerRadius = unit.Dp(12)
	btn.TextSize = unit.Sp(11)
	btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(4), Right: unit.Dp(4)}
	return btn.Layout(gtx)
}
