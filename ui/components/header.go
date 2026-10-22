package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type Header struct {
	MenuBtn    widget.Clickable
	RefreshBtn widget.Clickable
	ActionBtn  widget.Clickable
}

func (h *Header) Layout(
	gtx layout.Context,
	th *theme.Theme,
	title string,
	subtitle string,
	actionLabel string,
	onToggleMenu func(),
	onRefresh func(),
	onAction func(),
) layout.Dimensions {
	changed := false
	if h.MenuBtn.Clicked(gtx) && onToggleMenu != nil {
		onToggleMenu()
		changed = true
	}
	if h.RefreshBtn.Clicked(gtx) && onRefresh != nil {
		onRefresh()
		changed = true
	}
	if h.ActionBtn.Clicked(gtx) && onAction != nil {
		onAction()
		changed = true
	}
	if changed {
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(10),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceBetween,
		}.Layout(gtx,
			// Left: Menu Hamburger Button + Screen Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					// Hamburger Menu Button (44dp touch area)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &h.MenuBtn, "☰")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(10)
						btn.TextSize = unit.Sp(16)
						btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}
						return btn.Layout(gtx)
					}),

					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

					// Title & Optional Subtitle
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.H6(th.Theme, title)
								lbl.Color = theme.ColorTextPrimary
								lbl.TextSize = unit.Sp(17)
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if subtitle == "" {
									return layout.Dimensions{}
								}
								sub := material.Caption(th.Theme, subtitle)
								sub.Color = theme.ColorTextMuted
								sub.TextSize = unit.Sp(11)
								return sub.Layout(gtx)
							}),
						)
					}),
				)
			}),

			// Right: Action Buttons (e.g. Reload or Custom Action)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if onRefresh == nil {
							return layout.Dimensions{}
						}
						btn := material.Button(th.Theme, &h.RefreshBtn, "↻")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextMuted
						btn.CornerRadius = unit.Dp(10)
						btn.TextSize = unit.Sp(16)
						btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(11), Right: unit.Dp(11)}
						return btn.Layout(gtx)
					}),

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if actionLabel == "" || onAction == nil {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th.Theme, &h.ActionBtn, actionLabel)
							btn.Background = theme.ColorSurfaceElevated
							btn.Color = theme.ColorPrimary
							btn.CornerRadius = unit.Dp(10)
							btn.TextSize = unit.Sp(12)
							btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}
							return btn.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}
