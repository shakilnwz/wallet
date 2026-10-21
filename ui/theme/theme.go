package theme

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

//go:embed fonts/NotoSansBengali-Regular.ttf
var bengaliFontTTF []byte

var (
	ColorBackground      = color.NRGBA{R: 0x09, G: 0x0A, B: 0x0F, A: 0xFF}
	ColorSurface         = color.NRGBA{R: 0x12, G: 0x14, B: 0x1C, A: 0xFF}
	ColorSurfaceElevated = color.NRGBA{R: 0x1A, G: 0x1D, B: 0x29, A: 0xFF}
	ColorBorder          = color.NRGBA{R: 0x26, G: 0x2A, B: 0x3B, A: 0xFF}
	ColorTextPrimary     = color.NRGBA{R: 0xF4, G: 0xF4, B: 0xF5, A: 0xFF}
	ColorTextMuted       = color.NRGBA{R: 0xA1, G: 0xA1, B: 0xAA, A: 0xFF}
	ColorTextDim         = color.NRGBA{R: 0x71, G: 0x71, B: 0x7A, A: 0xFF}
	ColorIncome          = color.NRGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0xFF}
	ColorExpense         = color.NRGBA{R: 0xF4, G: 0x3F, B: 0x5E, A: 0xFF}
	ColorPrimary         = color.NRGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0xFF}
	ColorAmber           = color.NRGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}
)

type Theme struct {
	*material.Theme
}

func NewTheme() *Theme {
	var shaper *text.Shaper
	faces, err := opentype.ParseCollection(bengaliFontTTF)
	if err == nil {
		shaper = text.NewShaper(text.WithCollection(faces))
	} else {
		shaper = text.NewShaper()
	}

	th := material.NewTheme()
	th.Shaper = shaper
	th.Bg = ColorBackground
	th.Fg = ColorTextPrimary
	th.Palette.Fg = ColorTextPrimary
	th.Palette.Bg = ColorBackground
	th.Palette.ContrastBg = ColorPrimary
	th.Palette.ContrastFg = ColorTextPrimary
	return &Theme{Theme: th}
}

func FormatCurrency(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}
	intPart := int64(amount)
	decPart := int64(math.Round((amount - float64(intPart)) * 100))
	if decPart >= 100 {
		intPart++
		decPart = 0
	}

	str := fmt.Sprintf("%d", intPart)
	// Indian / Bangladeshi comma formatting (e.g. 1,00,000)
	var formatted strings.Builder
	n := len(str)
	if n <= 3 {
		formatted.WriteString(str)
	} else {
		// Last 3 digits
		last3 := str[n-3:]
		rest := str[:n-3]
		var parts []string
		for len(rest) > 2 {
			parts = append([]string{rest[len(rest)-2:]}, parts...)
			rest = rest[:len(rest)-2]
		}
		if len(rest) > 0 {
			parts = append([]string{rest}, parts...)
		}
		formatted.WriteString(strings.Join(parts, ","))
		formatted.WriteString(",")
		formatted.WriteString(last3)
	}

	return fmt.Sprintf("%s৳%s.%02d", sign, formatted.String(), decPart)
}

func FormatDateDisplay(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	entryDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if entryDay.Equal(today) {
		return "Today"
	}
	if entryDay.Equal(yesterday) {
		return "Yesterday"
	}
	return t.Format("02 Jan 2006")
}

func FormatMonthDisplay(monthStr string) string {
	if monthStr == "" || monthStr == "all" {
		return "All Time"
	}
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return monthStr
	}
	return t.Format("January 2006")
}

func Card(gtx layout.Context, bg color.NRGBA, radius unit.Dp, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rect := clip.RRect{
		Rect: image.Rectangle{Max: dims.Size},
		SE:   gtx.Dp(radius),
		SW:   gtx.Dp(radius),
		NW:   gtx.Dp(radius),
		NE:   gtx.Dp(radius),
	}

	paint.FillShape(gtx.Ops, bg, rect.Op(gtx.Ops))
	call.Add(gtx.Ops)
	return dims
}
