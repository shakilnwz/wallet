package components

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/ui/theme"
)

type TxList struct {
	List       layout.List
	RowBtns    []widget.Clickable
	DeleteBtns []widget.Clickable
}

func (l *TxList) Layout(
	gtx layout.Context,
	th *theme.Theme,
	transactions []models.Transaction,
	onEdit func(tx models.Transaction),
	onDelete func(tx models.Transaction),
) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	if len(transactions) == 0 {
		return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(24), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return theme.Card(gtx, theme.ColorSurface, unit.Dp(16), func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(28), Bottom: unit.Dp(28), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
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
		l.RowBtns = make([]widget.Clickable, len(transactions))
	}

	l.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return l.List.Layout(gtx, len(transactions), func(gtx layout.Context, idx int) layout.Dimensions {
			tx := transactions[idx]

			if l.DeleteBtns[idx].Clicked(gtx) && onDelete != nil {
				onDelete(tx)
				gtx.Execute(op.InvalidateCmd{})
			}
			if l.RowBtns[idx].Clicked(gtx) && onEdit != nil {
				onEdit(tx)
				gtx.Execute(op.InvalidateCmd{})
			}

			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return l.RowBtns[idx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
								// Left: Indicator icon + Title + Category + Date
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										// Circular indicator pill
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											dotColor := theme.ColorExpense
											dotSign := "-"
											if tx.Type == models.TypeIncome {
												dotColor = theme.ColorIncome
												dotSign = "+"
											}
											return theme.Card(gtx, dotColor, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													lbl := material.Caption(th.Theme, dotSign)
													lbl.Color = theme.ColorTextPrimary
													lbl.TextSize = unit.Sp(13)
													return lbl.Layout(gtx)
												})
											})
										}),

										layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

										// Text Details
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													title := material.Body1(th.Theme, tx.Benf)
													title.Color = theme.ColorTextPrimary
													title.TextSize = unit.Sp(15)
													return title.Layout(gtx)
												}),
												layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													subText := fmt.Sprintf("%s | %s | %s", tx.ExpCat, tx.Account, theme.FormatDateDisplay(tx.Date))
													sub := material.Caption(th.Theme, subText)
													sub.Color = theme.ColorTextMuted
													sub.TextSize = unit.Sp(11)
													return sub.Layout(gtx)
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													if tx.Notes == "" {
														return layout.Dimensions{}
													}
													note := material.Caption(th.Theme, tx.Notes)
													note.Color = theme.ColorTextDim
													note.TextSize = unit.Sp(11)
													note.MaxLines = 1
													return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, note.Layout)
												}),
											)
										}),
									)
								}),

								layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

								// Right: Amount + Vector Delete Icon Button
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
											val.TextSize = unit.Sp(14)
											return val.Layout(gtx)
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.IconButton(th.Theme, &l.DeleteBtns[idx], theme.IconClose, "Delete")
											btn.Background = theme.ColorSurfaceElevated
											btn.Color = theme.ColorExpense
											btn.Size = unit.Dp(14)
											btn.Inset = layout.UniformInset(unit.Dp(6))
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
	})
}
