package components

import (
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

type Keypad struct {
	DigitBtns [10]widget.Clickable // 0-9
	DotBtn    widget.Clickable
	BkspBtn   widget.Clickable

	QuickAddBtns [4]widget.Clickable // +100, +500, +1000, +5000
}

func (k *Keypad) Layout(gtx layout.Context, th *theme.Theme, currentVal string, onValChange func(string)) layout.Dimensions {
	// Handle digit clicks
	for i := 0; i <= 9; i++ {
		if k.DigitBtns[i].Clicked(gtx) {
			newVal := currentVal
			if newVal == "0" {
				newVal = strconv.Itoa(i)
			} else {
				newVal += strconv.Itoa(i)
			}
			onValChange(newVal)
		}
	}

	// Handle dot click
	if k.DotBtn.Clicked(gtx) {
		if !strings.Contains(currentVal, ".") {
			if currentVal == "" {
				onValChange("0.")
			} else {
				onValChange(currentVal + ".")
			}
		}
	}

	// Handle backspace click
	if k.BkspBtn.Clicked(gtx) {
		if len(currentVal) > 0 {
			onValChange(currentVal[:len(currentVal)-1])
		}
	}

	// Handle quick increment clicks
	quickAmts := []float64{100, 500, 1000, 5000}
	for i, amt := range quickAmts {
		if k.QuickAddBtns[i].Clicked(gtx) {
			cur, _ := strconv.ParseFloat(currentVal, 64)
			cur += amt
			onValChange(strconv.FormatFloat(cur, 'f', -1, 64))
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Quick Increment Chips
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					k.layoutQuickBtn(gtx, th, 0, "+100"),
					k.layoutQuickBtn(gtx, th, 1, "+500"),
					k.layoutQuickBtn(gtx, th, 2, "+1000"),
					k.layoutQuickBtn(gtx, th, 3, "+5000"),
				)
			})
		}),

		// Numeric Keypad Grid (Rows 1-4)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 1),
				k.layoutDigitBtn(gtx, th, 2),
				k.layoutDigitBtn(gtx, th, 3),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 4),
				k.layoutDigitBtn(gtx, th, 5),
				k.layoutDigitBtn(gtx, th, 6),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 7),
				k.layoutDigitBtn(gtx, th, 8),
				k.layoutDigitBtn(gtx, th, 9),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th.Theme, &k.DotBtn, ".")
					btn.Background = theme.ColorSurfaceElevated
					btn.Color = theme.ColorTextPrimary
					btn.TextSize = unit.Sp(18)
					btn.CornerRadius = unit.Dp(10)
					btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
				k.layoutDigitBtn(gtx, th, 0),
				layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th.Theme, &k.BkspBtn, "DEL")
					btn.Background = theme.ColorSurfaceElevated
					btn.Color = theme.ColorExpense
					btn.TextSize = unit.Sp(13)
					btn.CornerRadius = unit.Dp(10)
					btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}

func (k *Keypad) layoutDigitBtn(gtx layout.Context, th *theme.Theme, digit int) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th.Theme, &k.DigitBtns[digit], strconv.Itoa(digit))
		btn.Background = theme.ColorSurfaceElevated
		btn.Color = theme.ColorTextPrimary
		btn.TextSize = unit.Sp(18)
		btn.CornerRadius = unit.Dp(10)
		btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
		return btn.Layout(gtx)
	})
}

func (k *Keypad) layoutQuickBtn(gtx layout.Context, th *theme.Theme, idx int, text string) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th.Theme, &k.QuickAddBtns[idx], text)
		btn.Background = theme.ColorSurface
		btn.Color = theme.ColorTextMuted
		btn.TextSize = unit.Sp(11)
		btn.CornerRadius = unit.Dp(8)
		btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6)}
		return btn.Layout(gtx)
	})
}
