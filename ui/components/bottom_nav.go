package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
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
		Bottom: unit.Dp(10),
		Left:   unit.Dp(10),
		Right:  unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Stack{Alignment: layout.S}.Layout(gtx,
			// Bottom navigation surface.
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return theme.Card(gtx, theme.ColorSurface, unit.Dp(20), func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    unit.Dp(8),
						Bottom: unit.Dp(8),
						Left:   unit.Dp(8),
						Right:  unit.Dp(8),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
							Spacing:   layout.SpaceBetween,
							Gap:       4,
						}.Layout(gtx,
							// Home tab
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return b.layoutTabBtn(gtx, th, &b.HomeBtn, theme.IconHome, currentTab == TabHome)
							}),

							// Ledger tab
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return b.layoutTabBtn(gtx, th, &b.LedgerBtn, theme.IconLedger, currentTab == TabLedger)
							}),

							// Center gap for the floating add button.
							layout.Rigid(layout.Spacer{Width: unit.Dp(60)}.Layout),

							// Insights tab
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return b.layoutTabBtn(gtx, th, &b.InsightsBtn, theme.IconInsights, currentTab == TabInsights)
							}),

							// Menu tab (opens Drawer)
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return b.layoutTabBtn(gtx, th, &b.DrawerBtn, theme.IconMenu, currentTab == TabSync || currentTab == TabSettings)
							}),
						)
					})
				})
			}),

			// Floating add button. Its center sits on the nav surface's top edge,
			// so the button overlaps the bar instead of taking a slot in it.
			// Return zero dimensions so the overlay cannot affect the bar height;
			// the four tab buttons remain the bar's only measured children.
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				btn := material.IconButton(th.Theme, &b.AddBtn, theme.IconAdd, "Add Transaction")
				btn.Background = theme.ColorPrimary
				btn.Color = theme.ColorTextPrimary
				btn.Size = unit.Dp(36)
				btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(12), Right: unit.Dp(12)}

				macro := op.Record(gtx.Ops)
				dims := btn.Layout(gtx)
				call := macro.Stop()

				shift := op.Offset(image.Point{
					X: -dims.Size.X / 2,
					Y: -dims.Size.Y - gtx.Dp(unit.Dp(8)),
				}).Push(gtx.Ops)
				call.Add(gtx.Ops)
				shift.Pop()

				return layout.Dimensions{}
			}),
		)
	})
}

func (b *BottomNav) layoutTabBtn(
	gtx layout.Context,
	th *theme.Theme,
	click *widget.Clickable,
	ic *widget.Icon,
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
			return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(2), Right: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return theme.LayoutIcon(gtx, ic, col, unit.Dp(24))
						})
					}),
				)
			})
		})
	})
}
