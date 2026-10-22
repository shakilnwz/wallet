package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type MonthSelector struct {
	PrevBtn   widget.Clickable
	NextBtn   widget.Clickable
	ToggleBtn widget.Clickable
}

func (m *MonthSelector) Layout(
	gtx layout.Context,
	th *theme.Theme,
	selectedMonth string,
	availableMonths []string,
	onSelectMonth func(string),
) layout.Dimensions {
	changed := false

	if m.PrevBtn.Clicked(gtx) && len(availableMonths) > 0 {
		idx := -1
		for i, v := range availableMonths {
			if v == selectedMonth {
				idx = i
				break
			}
		}
		if idx != -1 && idx < len(availableMonths)-1 {
			if onSelectMonth != nil {
				onSelectMonth(availableMonths[idx+1])
				changed = true
			}
		}
	}

	if m.NextBtn.Clicked(gtx) && len(availableMonths) > 0 {
		idx := -1
		for i, v := range availableMonths {
			if v == selectedMonth {
				idx = i
				break
			}
		}
		if idx > 0 {
			if onSelectMonth != nil {
				onSelectMonth(availableMonths[idx-1])
				changed = true
			}
		}
	}

	if m.ToggleBtn.Clicked(gtx) {
		if selectedMonth == "all" {
			if len(availableMonths) > 0 {
				if onSelectMonth != nil {
					onSelectMonth(availableMonths[0])
					changed = true
				}
			}
		} else {
			if onSelectMonth != nil {
				onSelectMonth("all")
				changed = true
			}
		}
	}

	if changed {
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return theme.Card(gtx, theme.ColorSurface, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
				Left:   unit.Dp(6),
				Right:  unit.Dp(6),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					// Prev button (Vector Icon Button)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.IconButton(th.Theme, &m.PrevBtn, theme.IconChevronLeft, "Previous Month")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.Size = unit.Dp(18)
						btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(12), Right: unit.Dp(12)}
						return btn.Layout(gtx)
					}),

					// Center Month label (Tap to toggle All Time)
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Clickable(gtx, &m.ToggleBtn, func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Body1(th.Theme, theme.FormatMonthDisplay(selectedMonth))
									lbl.Color = theme.ColorTextPrimary
									lbl.TextSize = unit.Sp(14)
									return lbl.Layout(gtx)
								})
							})
						})
					}),

					// Next button (Vector Icon Button)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.IconButton(th.Theme, &m.NextBtn, theme.IconChevronRight, "Next Month")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.Size = unit.Dp(18)
						btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(12), Right: unit.Dp(12)}
						return btn.Layout(gtx)
					}),
				)
			})
		})
	})
}
