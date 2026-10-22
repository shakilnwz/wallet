package components

import (
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/money"
	"github.com/shakilnwz/wallet/ui/theme"
)

type Keypad struct {
	DigitBtns [10]widget.Clickable // 0-9
	DotBtn    widget.Clickable
	BkspBtn   widget.Clickable

	QuickAddBtns [4]widget.Clickable // +100, +500, +1000, +5000
}

func (k *Keypad) Layout(gtx layout.Context, th *theme.Theme, currentVal string, onValChange func(string)) layout.Dimensions {
	changed := false

	// Handle digit clicks
	for i := 0; i <= 9; i++ {
		if k.DigitBtns[i].Clicked(gtx) {
			// BDT has two decimal places (paisa); ignore further digits.
			if _, frac, ok := strings.Cut(currentVal, "."); ok && len(frac) >= 2 {
				continue
			}
			newVal := currentVal
			if newVal == "0" {
				newVal = strconv.Itoa(i)
			} else {
				newVal += strconv.Itoa(i)
			}
			if _, err := money.Parse(newVal); err != nil {
				continue // overflow
			}
			currentVal = newVal
			if onValChange != nil {
				onValChange(newVal)
			}
			changed = true
		}
	}

	// Handle dot click
	if k.DotBtn.Clicked(gtx) {
		if !strings.Contains(currentVal, ".") {
			var newVal string
			if currentVal == "" {
				newVal = "0."
			} else {
				newVal = currentVal + "."
			}
			currentVal = newVal
			if onValChange != nil {
				onValChange(newVal)
			}
			changed = true
		}
	}

	// Handle backspace click
	if k.BkspBtn.Clicked(gtx) {
		if len(currentVal) > 0 {
			newVal := currentVal[:len(currentVal)-1]
			currentVal = newVal
			if onValChange != nil {
				onValChange(newVal)
			}
			changed = true
		}
	}

	// Handle quick increment clicks
	quickAmts := []int64{100, 500, 1000, 5000}
	for i, amt := range quickAmts {
		if k.QuickAddBtns[i].Clicked(gtx) {
			cur, _ := money.Parse(currentVal) // empty or partial input counts as 0
			newVal := money.ToInput(cur + amt*100)
			currentVal = newVal
			if onValChange != nil {
				onValChange(newVal)
			}
			changed = true
		}
	}

	if changed {
		gtx.Execute(op.InvalidateCmd{})
	}

	const gapDp = unit.Dp(8)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Quick Increment Chips (4 equal columns with uniform spacing)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{Bottom: gapDp}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					k.layoutQuickBtn(gtx, th, 0, "+100"),
					layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
					k.layoutQuickBtn(gtx, th, 1, "+500"),
					layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
					k.layoutQuickBtn(gtx, th, 2, "+1000"),
					layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
					k.layoutQuickBtn(gtx, th, 3, "+5000"),
				)
			})
		}),

		// Numeric Keypad Grid: Row 1 (1, 2, 3)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 1),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 2),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 3),
			)
		}),
		layout.Rigid(layout.Spacer{Height: gapDp}.Layout),

		// Row 2 (4, 5, 6)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 4),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 5),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 6),
			)
		}),
		layout.Rigid(layout.Spacer{Height: gapDp}.Layout),

		// Row 3 (7, 8, 9)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				k.layoutDigitBtn(gtx, th, 7),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 8),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 9),
			)
		}),
		layout.Rigid(layout.Spacer{Height: gapDp}.Layout),

		// Row 4 (., 0, DEL) - all exact same flex 1 and exact same dimensions
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th.Theme, &k.DotBtn, ".")
					btn.Background = theme.ColorSurfaceElevated
					btn.Color = theme.ColorTextPrimary
					btn.TextSize = unit.Sp(20)
					btn.CornerRadius = unit.Dp(12)
					btn.Inset = layout.Inset{Top: unit.Dp(13), Bottom: unit.Dp(13)}
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				k.layoutDigitBtn(gtx, th, 0),
				layout.Rigid(layout.Spacer{Width: gapDp}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th.Theme, &k.BkspBtn, "DEL")
					btn.Background = theme.ColorSurfaceElevated
					btn.Color = theme.ColorExpense
					btn.TextSize = unit.Sp(13)
					btn.CornerRadius = unit.Dp(12)
					btn.Inset = layout.Inset{Top: unit.Dp(13), Bottom: unit.Dp(13)}
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
		btn.TextSize = unit.Sp(19)
		btn.CornerRadius = unit.Dp(12)
		btn.Inset = layout.Inset{Top: unit.Dp(13), Bottom: unit.Dp(13)}
		return btn.Layout(gtx)
	})
}

func (k *Keypad) layoutQuickBtn(gtx layout.Context, th *theme.Theme, idx int, text string) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(th.Theme, &k.QuickAddBtns[idx], text)
		btn.Background = theme.ColorSurfaceElevated
		btn.Color = theme.ColorTextMuted
		btn.TextSize = unit.Sp(11)
		btn.CornerRadius = unit.Dp(10)
		btn.Inset = layout.Inset{Top: unit.Dp(9), Bottom: unit.Dp(9)}
		return btn.Layout(gtx)
	})
}
