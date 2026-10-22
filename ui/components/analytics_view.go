package components

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/ui/theme"
)

type AnalyticsView struct {
	List layout.List
}

func (a *AnalyticsView) Layout(
	gtx layout.Context,
	th *theme.Theme,
	summary *models.MonthlySummary,
	overall *models.AccountSummary,
) layout.Dimensions {
	if summary == nil {
		summary = &models.MonthlySummary{}
	}
	if overall == nil {
		overall = &models.AccountSummary{}
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	a.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return a.List.Layout(gtx, 2, func(gtx layout.Context, idx int) layout.Dimensions {
			if idx == 0 {
				// Card 1: Category Spending Breakdown
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									title := material.H6(th.Theme, "Spending by Category")
									title.Color = theme.ColorTextPrimary
									title.TextSize = unit.Sp(16)
									return title.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if len(summary.Categories) == 0 {
										lbl := material.Caption(th.Theme, "No expense categories recorded for this month.")
										lbl.Color = theme.ColorTextMuted
										lbl.TextSize = unit.Sp(12)
										return lbl.Layout(gtx)
									}

									var rows []layout.FlexChild
									for _, cat := range summary.Categories {
										item := cat
										rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints.Min.X = gtx.Constraints.Max.X
												return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														gtx.Constraints.Min.X = gtx.Constraints.Max.X
														return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
															layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
																l := material.Body2(th.Theme, item.Category)
																l.Color = theme.ColorTextPrimary
																l.TextSize = unit.Sp(13)
																return l.Layout(gtx)
															}),
															layout.Rigid(func(gtx layout.Context) layout.Dimensions {
																valStr := fmt.Sprintf("%s (%.1f%%)", theme.FormatCurrency(item.Amount), item.Percentage)
																l := material.Body2(th.Theme, valStr)
																l.Color = theme.ColorTextMuted
																l.TextSize = unit.Sp(12)
																return l.Layout(gtx)
															}),
														)
													}),
													layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
													// Percentage progress bar (Full Width)
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														barWidth := gtx.Constraints.Max.X
														filledWidth := int(float64(barWidth) * (item.Percentage / 100.0))
														if filledWidth < 6 && item.Percentage > 0 {
															filledWidth = 6
														}
														h := gtx.Dp(unit.Dp(6))

														// Background rail
														rail := clip.RRect{
															Rect: image.Rect(0, 0, barWidth, h),
															SE:   h / 2, SW: h / 2, NW: h / 2, NE: h / 2,
														}
														paint.FillShape(gtx.Ops, theme.ColorSurfaceElevated, rail.Op(gtx.Ops))

														// Progress fill
														fill := clip.RRect{
															Rect: image.Rect(0, 0, filledWidth, h),
															SE:   h / 2, SW: h / 2, NW: h / 2, NE: h / 2,
														}
														paint.FillShape(gtx.Ops, theme.ColorExpense, fill.Op(gtx.Ops))

														return layout.Dimensions{Size: image.Point{X: barWidth, Y: h}}
													}),
												)
											})
										}))
									}
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
								}),
							)
						})
					})
				})
			}

			// Card 2: Account Holdings Distribution
			return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								title := material.H6(th.Theme, "Account Distribution")
								title.Color = theme.ColorTextPrimary
								title.TextSize = unit.Sp(16)
								return title.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														t := material.Caption(th.Theme, "Bank / MFS Account")
														t.Color = theme.ColorPrimary
														t.TextSize = unit.Sp(11)
														return t.Layout(gtx)
													}),
													layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														v := material.Body1(th.Theme, theme.FormatCurrency(overall.BankBalance))
														v.Color = theme.ColorTextPrimary
														v.TextSize = unit.Sp(15)
														return v.Layout(gtx)
													}),
												)
											})
										})
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														t := material.Caption(th.Theme, "Cash Holdings")
														t.Color = theme.ColorIncome
														t.TextSize = unit.Sp(11)
														return t.Layout(gtx)
													}),
													layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
													layout.Rigid(func(gtx layout.Context) layout.Dimensions {
														v := material.Body1(th.Theme, theme.FormatCurrency(overall.CashBalance))
														v.Color = theme.ColorTextPrimary
														v.TextSize = unit.Sp(15)
														return v.Layout(gtx)
													}),
												)
											})
										})
									}),
								)
							}),
						)
					})
				})
			})
		})
	})
}
