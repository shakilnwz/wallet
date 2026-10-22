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

type SettingsView struct {
	List layout.List

	ResetDataBtn   widget.Clickable
	CancelResetBtn widget.Clickable

	confirmReset bool
}

func (s *SettingsView) Layout(
	gtx layout.Context,
	th *theme.Theme,
	stats models.SyncStats,
	onResetData func(),
) layout.Dimensions {
	if s.ResetDataBtn.Clicked(gtx) && onResetData != nil {
		if s.confirmReset {
			s.confirmReset = false
			onResetData()
		} else {
			s.confirmReset = true
		}
		gtx.Execute(op.InvalidateCmd{})
	}
	if s.CancelResetBtn.Clicked(gtx) {
		s.confirmReset = false
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	s.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return s.List.Layout(gtx, 3, func(gtx layout.Context, idx int) layout.Dimensions {
			switch idx {
			case 0:
				// Card 1: Device Identity & DB Stats (Full Width)
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.H6(th.Theme, "Device Identity & Sync UUID")
									t.Color = theme.ColorTextPrimary
									t.TextSize = unit.Sp(16)
									return t.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											idLbl := material.Caption(th.Theme, stats.DeviceID)
											idLbl.Color = theme.ColorPrimary
											idLbl.TextSize = unit.Sp(12)
											return idLbl.Layout(gtx)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
								// Stats Grid: 3 equal cards
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Active")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(11)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.ActiveRecords))
															val.Color = theme.ColorIncome
															val.TextSize = unit.Sp(16)
															return val.Layout(gtx)
														}),
													)
												})
											})
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Tombstones")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(11)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.DeletedRecords))
															val.Color = theme.ColorTextMuted
															val.TextSize = unit.Sp(16)
															return val.Layout(gtx)
														}),
													)
												})
											})
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Total")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(11)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.TotalRecords))
															val.Color = theme.ColorPrimary
															val.TextSize = unit.Sp(16)
															return val.Layout(gtx)
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

			case 1:
				// Card 2: App & Architecture Info
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.H6(th.Theme, "About Wallet")
									t.Color = theme.ColorTextPrimary
									t.TextSize = unit.Sp(16)
									return t.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := "Wallet is a pure Go application built with Gio (gioui.org) and embedded SQLite (modernc.org/sqlite). All transactions are stored locally on your device in exact integer paisa without cloud intermediaries or external telemetry."
									lbl := material.Body2(th.Theme, txt)
									lbl.Color = theme.ColorTextMuted
									lbl.TextSize = unit.Sp(12)
									return lbl.Layout(gtx)
								}),
							)
						})
					})
				})

			case 2:
				// Card 3: Danger Zone: Delete All Transactions
				return layout.Inset{Bottom: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											title := material.H6(th.Theme, "Delete All Transactions")
											title.Color = theme.ColorExpense
											title.TextSize = unit.Sp(15)
											return title.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											if s.confirmReset {
												return layout.Dimensions{}
											}
											btn := material.Button(th.Theme, &s.ResetDataBtn, "Delete All")
											btn.Background = theme.ColorExpense
											btn.Color = theme.ColorTextPrimary
											btn.CornerRadius = unit.Dp(8)
											btn.TextSize = unit.Sp(11)
											btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}
											return btn.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := "Deletes every transaction on this device and writes tombstones that sync to paired devices."
									if s.confirmReset {
										txt = fmt.Sprintf("Delete all %d transactions? This action will sync to other paired devices.", stats.ActiveRecords)
									}
									d := material.Caption(th.Theme, txt)
									d.Color = theme.ColorTextMuted
									if s.confirmReset {
										d.Color = theme.ColorExpense
									}
									d.TextSize = unit.Sp(12)
									return d.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !s.confirmReset {
										return layout.Dimensions{}
									}
									return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												btn := material.Button(th.Theme, &s.CancelResetBtn, "Cancel")
												btn.Background = theme.ColorSurfaceElevated
												btn.Color = theme.ColorTextPrimary
												btn.CornerRadius = unit.Dp(10)
												btn.TextSize = unit.Sp(13)
												btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
												return btn.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												btn := material.Button(th.Theme, &s.ResetDataBtn, "Confirm Delete")
												btn.Background = theme.ColorExpense
												btn.Color = theme.ColorTextPrimary
												btn.CornerRadius = unit.Dp(10)
												btn.TextSize = unit.Sp(13)
												btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
												return btn.Layout(gtx)
											}),
										)
									})
								}),
							)
						})
					})
				})
			}
			return layout.Dimensions{}
		})
	})
}
