package components

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/ui/theme"
)

type TxList struct {
	List       layout.List
	DeleteBtns []widget.Clickable
}

func (l *TxList) Layout(
	gtx layout.Context,
	th *theme.Theme,
	transactions []models.Transaction,
	onDelete func(id string),
) layout.Dimensions {
	if len(transactions) == 0 {
		return layout.Inset{Top: unit.Dp(32), Bottom: unit.Dp(32), Left: unit.Dp(24), Right: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return theme.Card(gtx, theme.ColorSurface, unit.Dp(16), func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(24), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.Body1(th.Theme, "No Transactions Found")
							title.Color = theme.ColorTextPrimary
							title.TextSize = unit.Sp(15)
							return title.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							sub := material.Caption(th.Theme, "Tap the (+) button below to log your first transaction.")
							sub.Color = theme.ColorTextMuted
							sub.TextSize = unit.Sp(12)
							return sub.Layout(gtx)
						}),
					)
				})
			})
		})
	}

	// Ensure slice size
	if len(l.DeleteBtns) < len(transactions) {
		l.DeleteBtns = make([]widget.Clickable, len(transactions))
	}

	l.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return l.List.Layout(gtx, len(transactions), func(gtx layout.Context, idx int) layout.Dimensions {
			tx := transactions[idx]

			if l.DeleteBtns[idx].Clicked(gtx) && onDelete != nil {
				onDelete(tx.ID)
			}

			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return theme.Card(gtx, theme.ColorSurface, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
							// Left Column: Type indicator + Beneficiary + Category + Account
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									// Colored Circle Icon
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										dotColor := theme.ColorExpense
										dotSign := "↑"
										if tx.Type == models.TypeIncome {
											dotColor = theme.ColorIncome
											dotSign = "↓"
										}
										return theme.Card(gtx, dotColor, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												lbl := material.Caption(th.Theme, dotSign)
												lbl.Color = theme.ColorTextPrimary
												lbl.TextSize = unit.Sp(11)
												return lbl.Layout(gtx)
											})
										})
									}),

									layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),

									// Details
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												title := material.Body1(th.Theme, tx.Benf)
												title.Color = theme.ColorTextPrimary
												title.TextSize = unit.Sp(14)
												return title.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												subText := fmt.Sprintf("%s • %s • %s", tx.ExpCat, tx.Account, theme.FormatDateDisplay(tx.Date))
												sub := material.Caption(th.Theme, subText)
												sub.Color = theme.ColorTextMuted
												sub.TextSize = unit.Sp(10)
												return sub.Layout(gtx)
											}),
										)
									}),
								)
							}),

							// Right Column: Formatted Amount + Delete Button
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										amtColor := theme.ColorTextPrimary
										prefix := "-"
										if tx.Type == models.TypeIncome {
											amtColor = theme.ColorIncome
											prefix = "+"
										}
										amtStr := prefix + theme.FormatCurrency(tx.Amount)
										val := material.Body1(th.Theme, amtStr)
										val.Color = amtColor
										val.TextSize = unit.Sp(13)
										return val.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btn := material.Button(th.Theme, &l.DeleteBtns[idx], "✕")
										btn.Background = theme.ColorSurfaceElevated
										btn.Color = theme.ColorExpense
										btn.CornerRadius = unit.Dp(8)
										btn.TextSize = unit.Sp(11)
										btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(7), Right: unit.Dp(7)}
										return btn.Layout(gtx)
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
