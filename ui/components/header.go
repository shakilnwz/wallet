package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type Header struct {
	RefreshBtn widget.Clickable
	SyncBtn    widget.Clickable
}

func (h *Header) Layout(gtx layout.Context, th *theme.Theme, onRefresh func(), onOpenSync func()) layout.Dimensions {
	if h.RefreshBtn.Clicked(gtx) && onRefresh != nil {
		onRefresh()
	}
	if h.SyncBtn.Clicked(gtx) && onOpenSync != nil {
		onOpenSync()
	}

	return layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceBetween,
		}.Layout(gtx,
			// Title & SQLite Badge
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.H6(th.Theme, "WALLET")
						lbl.Color = theme.ColorTextPrimary
						lbl.TextSize = unit.Sp(18)
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(2),
								Bottom: unit.Dp(2),
								Left:   unit.Dp(6),
								Right:  unit.Dp(6),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								badge := material.Caption(th.Theme, "SQLite Go")
								badge.Color = theme.ColorIncome
								badge.TextSize = unit.Sp(10)
								return badge.Layout(gtx)
							})
						})
					}),
				)
			}),

			// Action Buttons (Refresh & Sync)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &h.RefreshBtn, "Reload")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(8)
						btn.TextSize = unit.Sp(12)
						btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &h.SyncBtn, "Sync")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorPrimary
						btn.CornerRadius = unit.Dp(8)
						btn.TextSize = unit.Sp(12)
						btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
}
