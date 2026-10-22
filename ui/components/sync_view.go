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

type SyncView struct {
	List layout.List

	PeerAddrEditor widget.Editor
	SyncNowBtn     widget.Clickable
	ToggleServer   widget.Clickable
	ResetDataBtn   widget.Clickable

	StatusMsg string
}

func (s *SyncView) Layout(
	gtx layout.Context,
	th *theme.Theme,
	stats models.SyncStats,
	serverRunning bool,
	serverPort int,
	onToggleServer func(),
	onSyncWithPeer func(peerAddr string),
	onResetData func(),
) layout.Dimensions {
	if s.ToggleServer.Clicked(gtx) && onToggleServer != nil {
		onToggleServer()
	}

	if s.SyncNowBtn.Clicked(gtx) && onSyncWithPeer != nil {
		addr := s.PeerAddrEditor.Text()
		if addr != "" {
			onSyncWithPeer(addr)
		}
	}

	if s.ResetDataBtn.Clicked(gtx) && onResetData != nil {
		onResetData()
	}

	s.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return s.List.Layout(gtx, 4, func(gtx layout.Context, idx int) layout.Dimensions {
			switch idx {
			case 0:
				// Card 1: Device Identity & DB Stats
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.H6(th.Theme, "Device Identity & Sync UUID")
									t.Color = theme.ColorTextPrimary
									t.TextSize = unit.Sp(15)
									return t.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											idLbl := material.Caption(th.Theme, stats.DeviceID)
											idLbl.Color = theme.ColorPrimary
											idLbl.TextSize = unit.Sp(11)
											return idLbl.Layout(gtx)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
								// Stats Grid
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Active")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(10)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.ActiveRecords))
															val.Color = theme.ColorIncome
															val.TextSize = unit.Sp(14)
															return val.Layout(gtx)
														}),
													)
												})
											})
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Tombstones")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(10)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.DeletedRecords))
															val.Color = theme.ColorTextMuted
															val.TextSize = unit.Sp(14)
															return val.Layout(gtx)
														}),
													)
												})
											})
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															lbl := material.Caption(th.Theme, "Total")
															lbl.Color = theme.ColorTextMuted
															lbl.TextSize = unit.Sp(10)
															return lbl.Layout(gtx)
														}),
														layout.Rigid(func(gtx layout.Context) layout.Dimensions {
															val := material.Body1(th.Theme, fmt.Sprintf("%d", stats.TotalRecords))
															val.Color = theme.ColorPrimary
															val.TextSize = unit.Sp(14)
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
				// Card 2: Local Network P2P Sync Server
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											title := material.H6(th.Theme, "Local Sync Server")
											title.Color = theme.ColorTextPrimary
											title.TextSize = unit.Sp(15)
											return title.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btnText := "Start Server"
											btnBg := theme.ColorPrimary
											if serverRunning {
												btnText = fmt.Sprintf("Stop (: %d)", serverPort)
												btnBg = theme.ColorIncome
											}
											btn := material.Button(th.Theme, &s.ToggleServer, btnText)
											btn.Background = btnBg
											btn.Color = theme.ColorTextPrimary
											btn.CornerRadius = unit.Dp(8)
											btn.TextSize = unit.Sp(11)
											btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}
											return btn.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := "Start the local sync server to allow other devices on your Wi-Fi/LAN to discover and merge records with this device."
									if serverRunning {
										txt = fmt.Sprintf("Server running on port %d. Other devices can sync by entering this device's local IP and port.", serverPort)
									}
									desc := material.Caption(th.Theme, txt)
									desc.Color = theme.ColorTextMuted
									desc.TextSize = unit.Sp(11)
									return desc.Layout(gtx)
								}),
							)
						})
					})
				})

			case 2:
				// Card 3: Connect & Sync with Remote Peer
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.H6(th.Theme, "Sync With Peer")
									t.Color = theme.ColorTextPrimary
									t.TextSize = unit.Sp(15)
									return t.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								// Peer Address Input
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											ed := material.Editor(th.Theme, &s.PeerAddrEditor, "e.g. 192.168.1.100:8080")
											ed.Color = theme.ColorTextPrimary
											ed.TextSize = unit.Sp(13)
											return ed.Layout(gtx)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th.Theme, &s.SyncNowBtn, "Two-Way Sync Now")
									btn.Background = theme.ColorPrimary
									btn.Color = theme.ColorTextPrimary
									btn.CornerRadius = unit.Dp(10)
									btn.TextSize = unit.Sp(13)
									btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}
									return btn.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if s.StatusMsg == "" {
										return layout.Dimensions{}
									}
									status := material.Caption(th.Theme, s.StatusMsg)
									status.Color = theme.ColorIncome
									status.TextSize = unit.Sp(11)
									return status.Layout(gtx)
								}),
							)
						})
					})
				})

			case 3:
				// Card 4: Danger Zone: Reset Data
				return layout.Inset{Bottom: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											title := material.H6(th.Theme, "Delete All Transactions")
											title.Color = theme.ColorExpense
											title.TextSize = unit.Sp(14)
											return title.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.Button(th.Theme, &s.ResetDataBtn, "Delete All")
											btn.Background = theme.ColorExpense
											btn.Color = theme.ColorTextPrimary
											btn.CornerRadius = unit.Dp(8)
											btn.TextSize = unit.Sp(11)
											btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}
											return btn.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									d := material.Caption(th.Theme, "Deletes every transaction here, and on paired devices at their next sync.")
									d.Color = theme.ColorTextMuted
									d.TextSize = unit.Sp(11)
									return d.Layout(gtx)
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
