package components

import (
	"gioui.org/layout"
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
	if m.PrevBtn.Clicked(gtx) && len(availableMonths) > 0 {
		idx := -1
		for i, v := range availableMonths {
			if v == selectedMonth {
				idx = i
				break
			}
		}
		if idx != -1 && idx < len(availableMonths)-1 {
			onSelectMonth(availableMonths[idx+1])
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
			onSelectMonth(availableMonths[idx-1])
		}
	}

	if m.ToggleBtn.Clicked(gtx) {
		if selectedMonth == "all" {
			if len(availableMonths) > 0 {
				onSelectMonth(availableMonths[0])
			}
		} else {
			onSelectMonth("all")
		}
	}

	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return theme.Card(gtx, theme.ColorSurface, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(6),
				Bottom: unit.Dp(6),
				Left:   unit.Dp(8),
				Right:  unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					// Prev button
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &m.PrevBtn, "‹")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(8)
						btn.TextSize = unit.Sp(18)
						btn.Inset = layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(10), Right: unit.Dp(10)}
						return btn.Layout(gtx)
					}),

					// Center Month label
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Clickable(gtx, &m.ToggleBtn, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(th.Theme, theme.FormatMonthDisplay(selectedMonth))
							lbl.Color = theme.ColorTextPrimary
							lbl.TextSize = unit.Sp(14)
							return lbl.Layout(gtx)
						})
						return btn
					}),

					// Next button
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th.Theme, &m.NextBtn, "›")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextPrimary
						btn.CornerRadius = unit.Dp(8)
						btn.TextSize = unit.Sp(18)
						btn.Inset = layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(10), Right: unit.Dp(10)}
						return btn.Layout(gtx)
					}),
				)
			})
		})
	})
}
