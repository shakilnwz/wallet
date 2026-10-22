package components

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/ui/theme"
)

type HeroCard struct{}

func (h *HeroCard) Layout(gtx layout.Context, th *theme.Theme, summary *models.AccountSummary, monthLabel string) layout.Dimensions {
	if summary == nil {
		summary = &models.AccountSummary{}
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(14),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return theme.Card(gtx, theme.ColorSurface, unit.Dp(20), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(18),
				Bottom: unit.Dp(18),
				Left:   unit.Dp(18),
				Right:  unit.Dp(18),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Header row: Net Balance title + Balance Status
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								caption := material.Caption(th.Theme, fmt.Sprintf("NET BALANCE (%s)", monthLabel))
								caption.Color = theme.ColorTextMuted
								caption.TextSize = unit.Sp(11)
								return caption.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								statusText := "balanced [ok]"
								statusColor := theme.ColorIncome
								if !summary.IsBalanced {
									statusText = fmt.Sprintf("reconcile %s", theme.FormatCurrency(summary.NeedBalancing))
									statusColor = theme.ColorAmber
								}
								lbl := material.Caption(th.Theme, statusText)
								lbl.Color = statusColor
								lbl.TextSize = unit.Sp(11)
								return lbl.Layout(gtx)
							}),
						)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

					// Big Net Balance Figure
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						val := material.H5(th.Theme, theme.FormatCurrency(summary.TotalBalance))
						val.Color = theme.ColorTextPrimary
						val.TextSize = unit.Sp(28)
						return val.Layout(gtx)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),

					// Income vs Expense Split Row (Full Width)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							// Income Column
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												t := material.Caption(th.Theme, "+ Income")
												t.Color = theme.ColorIncome
												t.TextSize = unit.Sp(11)
												return t.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												v := material.Body1(th.Theme, theme.FormatCurrency(summary.TotalIncome))
												v.Color = theme.ColorIncome
												v.TextSize = unit.Sp(14)
												return v.Layout(gtx)
											}),
										)
									})
								})
							}),

							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),

							// Expense Column
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												t := material.Caption(th.Theme, "- Expense")
												t.Color = theme.ColorExpense
												t.TextSize = unit.Sp(11)
												return t.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												v := material.Body1(th.Theme, theme.FormatCurrency(summary.TotalExpense))
												v.Color = theme.ColorExpense
												v.TextSize = unit.Sp(14)
												return v.Layout(gtx)
											}),
										)
									})
								})
							}),
						)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),

					// Account vs Cash Distribution Footer (Full Width)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								txt := fmt.Sprintf("Bank / MFS: %s", theme.FormatCurrency(summary.BankBalance))
								lbl := material.Caption(th.Theme, txt)
								lbl.Color = theme.ColorPrimary
								lbl.TextSize = unit.Sp(11)
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								txt := fmt.Sprintf("Cash: %s", theme.FormatCurrency(summary.CashBalance))
								lbl := material.Caption(th.Theme, txt)
								lbl.Color = theme.ColorTextMuted
								lbl.TextSize = unit.Sp(11)
								return lbl.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
	})
}
