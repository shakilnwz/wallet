package theme

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func mustIcon(data []byte) *widget.Icon {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return ic
}

var (
	IconMenu         = mustIcon(icons.NavigationMenu)
	IconRefresh      = mustIcon(icons.NavigationRefresh)
	IconClose        = mustIcon(icons.NavigationClose)
	IconChevronLeft  = mustIcon(icons.NavigationChevronLeft)
	IconChevronRight = mustIcon(icons.NavigationChevronRight)
	IconAdd          = mustIcon(icons.ContentAdd)
	IconDelete       = mustIcon(icons.ActionDelete)
	IconHome         = mustIcon(icons.ActionHome)
	IconLedger       = mustIcon(icons.ActionDescription)
	IconInsights     = mustIcon(icons.ActionAssessment)
	IconSync         = mustIcon(icons.NotificationSync)
	IconSettings     = mustIcon(icons.ActionSettings)
)

func LayoutIcon(gtx layout.Context, ic *widget.Icon, col color.NRGBA, size unit.Dp) layout.Dimensions {
	px := gtx.Dp(size)
	gtx.Constraints.Min = image.Pt(px, px)
	gtx.Constraints.Max = image.Pt(px, px)
	return ic.Layout(gtx, col)
}
