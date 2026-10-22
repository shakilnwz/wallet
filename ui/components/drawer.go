package components

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type DrawerItem int

const (
	DrawerItemHome DrawerItem = iota
	DrawerItemLedger
	DrawerItemInsights
	DrawerItemSync
	DrawerItemSettings
)

type NavDrawer struct {
	IsOpen bool

	HomeBtn     widget.Clickable
	LedgerBtn   widget.Clickable
	InsightsBtn widget.Clickable
	SyncBtn     widget.Clickable
	SettingsBtn widget.Clickable
	CloseBtn    widget.Clickable

	scrim  widget.Clickable
	drawer widget.Clickable

	anim *FloatAnim
}

func (d *NavDrawer) ensureAnim() {
	if d.anim == nil {
		d.anim = NewFloatAnim(0, 220*time.Millisecond)
	}
}

func (d *NavDrawer) Open() {
	d.IsOpen = true
	d.ensureAnim()
	d.anim.SetTarget(1.0)
}

func (d *NavDrawer) Close() {
	d.IsOpen = false
	d.ensureAnim()
	d.anim.SetTarget(0.0)
}

func (d *NavDrawer) Toggle() {
	if d.IsOpen {
		d.Close()
	} else {
		d.Open()
	}
}

func (d *NavDrawer) Layout(
	gtx layout.Context,
	th *theme.Theme,
	currentTab Tab,
	localIP string,
	totalTxCount int,
	onSelectTab func(Tab),
) layout.Dimensions {
	d.ensureAnim()
	progress := d.anim.Update(gtx)

	if !d.IsOpen && progress <= 0.001 {
		return layout.Dimensions{}
	}

	// Handle Back button to close drawer
	if d.IsOpen {
		for {
			ev, ok := gtx.Event(key.Filter{Name: key.NameBack}, key.Filter{Name: key.NameEscape})
			if !ok {
				break
			}
			if e, ok := ev.(key.Event); ok && e.State == key.Press {
				d.Close()
				gtx.Execute(op.InvalidateCmd{})
			}
		}
	}

	// Scrim and close button clicks
	if d.scrim.Clicked(gtx) || d.CloseBtn.Clicked(gtx) {
		d.Close()
		gtx.Execute(op.InvalidateCmd{})
	}

	// Navigation item clicks
	selectItem := func(t Tab) {
		d.Close()
		gtx.Execute(op.InvalidateCmd{})
		if onSelectTab != nil {
			onSelectTab(t)
		}
	}

	if d.HomeBtn.Clicked(gtx) {
		selectItem(TabHome)
	}
	if d.LedgerBtn.Clicked(gtx) {
		selectItem(TabLedger)
	}
	if d.InsightsBtn.Clicked(gtx) {
		selectItem(TabInsights)
	}
	if d.SyncBtn.Clicked(gtx) {
		selectItem(TabSync)
	}
	if d.SettingsBtn.Clicked(gtx) {
		selectItem(TabSettings)
	}

	easedProgress := EaseOutCubic(progress)

	// Dimmed backdrop scrim
	scrimAlpha := uint8(float32(0xB0) * easedProgress)
	if scrimAlpha > 0 {
		scrimPaint := color.NRGBA{A: scrimAlpha}
		gtx.Constraints.Min = gtx.Constraints.Max
		d.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			paint.FillShape(gtx.Ops, scrimPaint, clip.Rect{Max: size}.Op())
			return layout.Dimensions{Size: size}
		})
	}

	// Drawer container width: 80% of screen width, max 320dp
	drawerWidth := gtx.Constraints.Max.X * 80 / 100
	maxDrawerPx := gtx.Dp(unit.Dp(320))
	if drawerWidth > maxDrawerPx {
		drawerWidth = maxDrawerPx
	}

	offsetX := int(float32(drawerWidth) * (easedProgress - 1.0))

	return layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = drawerWidth
		gtx.Constraints.Max.X = drawerWidth
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

		trans := op.Offset(image.Pt(offsetX, 0)).Push(gtx.Ops)
		defer trans.Pop()

		return d.drawer.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return theme.Card(gtx, theme.ColorSurface, 0, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Drawer Header
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return d.layoutHeader(gtx, th, totalTxCount)
					}),

					// Subtle divider line
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lineHeight := gtx.Dp(unit.Dp(1))
						paint.FillShape(gtx.Ops, theme.ColorBorder, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, lineHeight)}.Op())
						return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, lineHeight)}
					}),

					// Navigation Menu Items
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return d.layoutMenu(gtx, th, currentTab)
					}),

					// Drawer Footer
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return d.layoutFooter(gtx, th, localIP)
					}),
				)
			})
		})
	})
}

func (d *NavDrawer) layoutHeader(gtx layout.Context, th *theme.Theme, totalTxCount int) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(24),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// App Logo icon card
								return theme.Card(gtx, theme.ColorPrimary, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
									return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										icon := material.Body1(th.Theme, "৳")
										icon.Color = theme.ColorTextPrimary
										icon.TextSize = unit.Sp(16)
										return icon.Layout(gtx)
									})
								})
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										title := material.H6(th.Theme, "WALLET")
										title.Color = theme.ColorTextPrimary
										title.TextSize = unit.Sp(18)
										return title.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										sub := material.Caption(th.Theme, "Offline-First Ledger")
										sub.Color = theme.ColorTextMuted
										sub.TextSize = unit.Sp(11)
										return sub.Layout(gtx)
									}),
								)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.IconButton(th.Theme, &d.CloseBtn, theme.IconClose, "Close")
						btn.Background = theme.ColorSurfaceElevated
						btn.Color = theme.ColorTextMuted
						btn.Size = unit.Dp(16)
						btn.Inset = layout.UniformInset(unit.Dp(8))
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						txt := fmt.Sprintf("SQLite Active Records: %d", totalTxCount)
						lbl := material.Caption(th.Theme, txt)
						lbl.Color = theme.ColorIncome
						lbl.TextSize = unit.Sp(11)
						return lbl.Layout(gtx)
					})
				})
			}),
		)
	})
}

func (d *NavDrawer) layoutMenu(gtx layout.Context, th *theme.Theme, currentTab Tab) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(16),
		Bottom: unit.Dp(16),
		Left:   unit.Dp(12),
		Right:  unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.layoutMenuItem(gtx, th, &d.HomeBtn, theme.IconHome, "Dashboard", "Overview & Recent", currentTab == TabHome)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.layoutMenuItem(gtx, th, &d.LedgerBtn, theme.IconLedger, "Ledger", "All Transactions & History", currentTab == TabLedger)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.layoutMenuItem(gtx, th, &d.InsightsBtn, theme.IconInsights, "Insights", "Spending Analytics & Charts", currentTab == TabInsights)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.layoutMenuItem(gtx, th, &d.SyncBtn, theme.IconSync, "P2P Sync", "Local Wi-Fi Peer Sync", currentTab == TabSync)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.layoutMenuItem(gtx, th, &d.SettingsBtn, theme.IconSettings, "Data & Settings", "DB Stats, Backup & Reset", currentTab == TabSettings)
			}),
		)
	})
}

func (d *NavDrawer) layoutMenuItem(
	gtx layout.Context,
	th *theme.Theme,
	click *widget.Clickable,
	ic *widget.Icon,
	title string,
	subtitle string,
	isActive bool,
) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	bg := theme.ColorSurface
	titleColor := theme.ColorTextMuted
	subColor := theme.ColorTextDim
	iconColor := theme.ColorTextMuted

	if isActive {
		bg = theme.ColorSurfaceElevated
		titleColor = theme.ColorTextPrimary
		subColor = theme.ColorPrimary
		iconColor = theme.ColorPrimary
	}

	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return theme.Card(gtx, bg, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(12),
				Bottom: unit.Dp(12),
				Left:   unit.Dp(14),
				Right:  unit.Dp(14),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return theme.LayoutIcon(gtx, ic, iconColor, unit.Dp(20))
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								t := material.Body1(th.Theme, title)
								t.Color = titleColor
								t.TextSize = unit.Sp(15)
								return t.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								s := material.Caption(th.Theme, subtitle)
								s.Color = subColor
								s.TextSize = unit.Sp(11)
								return s.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
	})
}

func (d *NavDrawer) layoutFooter(gtx layout.Context, th *theme.Theme, localIP string) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
		Right:  unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		ipText := "Wi-Fi: Not connected"
		if localIP != "" {
			ipText = "Wi-Fi IP: " + localIP
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th.Theme, ipText)
				lbl.Color = theme.ColorTextDim
				lbl.TextSize = unit.Sp(11)
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th.Theme, "Wallet v0.2.0 | Native Go")
				lbl.Color = theme.ColorTextDim
				lbl.TextSize = unit.Sp(10)
				return lbl.Layout(gtx)
			}),
		)
	})
}
