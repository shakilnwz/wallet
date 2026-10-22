package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
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
	TabSettings
)

type BottomNav struct {
	HomeBtn     widget.Clickable
	LedgerBtn   widget.Clickable
	InsightsBtn widget.Clickable
	DrawerBtn   widget.Clickable
	AddBtn      widget.Clickable
}

func (b *BottomNav) Layout(
	gtx layout.Context,
	th *theme.Theme,
	currentTab Tab,
	onSelectTab func(Tab),
	onOpenAdd func(),
	onToggleDrawer func(),
) layout.Dimensions {
	changed := false
	if b.HomeBtn.Clicked(gtx) && onSelectTab != nil {
		onSelectTab(TabHome)
		changed = true
	}
	if b.LedgerBtn.Clicked(gtx) && onSelectTab != nil {
		onSelectTab(TabLedger)
		changed = true
	}
	if b.InsightsBtn.Clicked(gtx) && onSelectTab != nil {
		onSelectTab(TabInsights)
		changed = true
	}
	if b.DrawerBtn.Clicked(gtx) && onToggleDrawer != nil {
		onToggleDrawer()
		changed = true
	}
	if b.AddBtn.Clicked(gtx) && onOpenAdd != nil {
		onOpenAdd()
		changed = true
	}

	if changed {
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Inset{
		Bottom: unit.Dp(12),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return theme.Card(gtx, theme.ColorSurface, unit.Dp(24), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(6),
				Bottom: unit.Dp(6),
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					// Home tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.HomeBtn, theme.IconHome, "Home", currentTab == TabHome)
					}),

					// Ledger tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.LedgerBtn, theme.IconLedger, "Ledger", currentTab == TabLedger)
					}),

					// Center Prominent Add (+) Button (Large touch target with vector icon)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(6), Right: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := material.IconButton(th.Theme, &b.AddBtn, theme.IconAdd, "Add Transaction")
							btn.Background = theme.ColorPrimary
							btn.Color = theme.ColorTextPrimary
							btn.Size = unit.Dp(22)
							btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(18), Right: unit.Dp(18)}
							return btn.Layout(gtx)
						})
					}),

					// Insights tab
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.InsightsBtn, theme.IconInsights, "Insights", currentTab == TabInsights)
					}),

					// Menu tab (opens Drawer)
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return b.layoutTabBtn(gtx, th, &b.DrawerBtn, theme.IconMenu, "Menu", currentTab == TabSync || currentTab == TabSettings)
					}),
				)
			})
		})
	})
}

func (b *BottomNav) layoutTabBtn(
	gtx layout.Context,
	th *theme.Theme,
	click *widget.Clickable,
	ic *widget.Icon,
	title string,
	isActive bool,
) layout.Dimensions {
	col := theme.ColorTextMuted
	bg := theme.ColorSurface
	if isActive {
		col = theme.ColorTextPrimary
		bg = theme.ColorSurfaceElevated
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return theme.Card(gtx, bg, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(2), Right: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return theme.LayoutIcon(gtx, ic, col, unit.Dp(20))
						})
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(th.Theme, title)
							lbl.Color = col
							lbl.TextSize = unit.Sp(10)
							lbl.Alignment = text.Middle
							return lbl.Layout(gtx)
						})
					}),
				)
			})
		})
	})
}
