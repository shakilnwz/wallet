package components

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type SyncView struct {
	List layout.List

	PeerAddrEditor widget.Editor
	SyncNowBtn     widget.Clickable
	ToggleServer   widget.Clickable

	StatusMsg string
	StatusErr bool
	Syncing   bool
}

// SetStatus shows a sync status line, in red when isErr.
func (s *SyncView) SetStatus(msg string, isErr bool) {
	s.StatusMsg, s.StatusErr = msg, isErr
}

func (s *SyncView) Layout(
	gtx layout.Context,
	th *theme.Theme,
	localIP string,
	serverRunning bool,
	serverPort int,
	onToggleServer func(),
	onSyncWithPeer func(peerAddr string),
) layout.Dimensions {
	if s.ToggleServer.Clicked(gtx) && onToggleServer != nil {
		onToggleServer()
		gtx.Execute(op.InvalidateCmd{})
	}

	s.PeerAddrEditor.SingleLine = true
	if s.SyncNowBtn.Clicked(gtx) && onSyncWithPeer != nil && !s.Syncing {
		addr := strings.TrimSpace(s.PeerAddrEditor.Text())
		if addr == "" {
			s.SetStatus("Enter the other device's address first.", true)
		} else {
			onSyncWithPeer(addr)
		}
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	s.List.Axis = layout.Vertical

	return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return s.List.Layout(gtx, 2, func(gtx layout.Context, idx int) layout.Dimensions {
			switch idx {
			case 0:
				// Card 1: Local Network P2P Sync Server (Full Width)
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											title := material.H6(th.Theme, "Local Sync Server")
											title.Color = theme.ColorTextPrimary
											title.TextSize = unit.Sp(16)
											return title.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btnText := "Start Server"
											btnBg := theme.ColorPrimary
											if serverRunning {
												btnText = "Stop Server"
												btnBg = theme.ColorExpense
											}
											btn := material.Button(th.Theme, &s.ToggleServer, btnText)
											btn.Background = btnBg
											btn.Color = theme.ColorTextPrimary
											btn.CornerRadius = unit.Dp(10)
											btn.TextSize = unit.Sp(12)
											btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(14), Right: unit.Dp(14)}
											return btn.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									txt := "Start the server so other devices on the same Wi-Fi can sync with this one."
									if serverRunning && localIP == "" {
										txt = fmt.Sprintf("Running on port %d, but this device doesn't seem to be on a Wi-Fi network.", serverPort)
									} else if serverRunning {
										txt = "On the other device, enter this address and tap Sync:"
									}
									desc := material.Caption(th.Theme, txt)
									desc.Color = theme.ColorTextMuted
									desc.TextSize = unit.Sp(12)
									return desc.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !serverRunning || localIP == "" {
										return layout.Dimensions{}
									}
									return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
											return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints.Min.X = gtx.Constraints.Max.X
												lbl := material.H6(th.Theme, fmt.Sprintf("%s:%d", localIP, serverPort))
												lbl.Color = theme.ColorPrimary
												lbl.TextSize = unit.Sp(18)
												return lbl.Layout(gtx)
											})
										})
									})
								}),
							)
						})
					})
				})

			case 1:
				// Card 2: Connect & Sync with Remote Peer (Full Width)
				return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return theme.Card(gtx, theme.ColorSurface, unit.Dp(18), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									t := material.H6(th.Theme, "Sync With Peer")
									t.Color = theme.ColorTextPrimary
									t.TextSize = unit.Sp(16)
									return t.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									desc := material.Caption(th.Theme, "Enter peer IP and port from the other device running Wallet:")
									desc.Color = theme.ColorTextMuted
									desc.TextSize = unit.Sp(12)
									return desc.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								// Peer Address Input (Full Width)
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(14), Right: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											ed := material.Editor(th.Theme, &s.PeerAddrEditor, "e.g. 192.168.1.100:8080")
											ed.Color = theme.ColorTextPrimary
											ed.TextSize = unit.Sp(14)
											return ed.Layout(gtx)
										})
									})
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label, bg := "Two-Way Sync Now", theme.ColorPrimary
									if s.Syncing {
										label, bg = "Syncing...", theme.ColorSurfaceElevated
									}
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									btn := material.Button(th.Theme, &s.SyncNowBtn, label)
									btn.Background = bg
									btn.Color = theme.ColorTextPrimary
									btn.CornerRadius = unit.Dp(12)
									btn.TextSize = unit.Sp(14)
									btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
									return btn.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if s.StatusMsg == "" {
										return layout.Dimensions{}
									}
									status := material.Caption(th.Theme, s.StatusMsg)
									status.Color = theme.ColorIncome
									if s.StatusErr {
										status.Color = theme.ColorExpense
									} else if s.Syncing {
										status.Color = theme.ColorTextMuted
									}
									status.TextSize = unit.Sp(12)
									return status.Layout(gtx)
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
