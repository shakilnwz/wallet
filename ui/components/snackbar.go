package components

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/ui/theme"
)

// Snackbar shows a short message at the bottom of the screen, optionally with
// one action (e.g. Undo). It hides itself after its duration.
type Snackbar struct {
	msg      string
	isError  bool
	action   string
	onAction func()
	until    time.Time
	btn      widget.Clickable
}

// Show displays msg for d. action and onAction may be empty/nil.
func (s *Snackbar) Show(msg string, isError bool, d time.Duration, action string, onAction func()) {
	s.msg, s.isError, s.action, s.onAction = msg, isError, action, onAction
	s.until = time.Now().Add(d)
}

func (s *Snackbar) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	if s.msg == "" {
		return layout.Dimensions{}
	}
	if !gtx.Now.Before(s.until) {
		s.msg = ""
		return layout.Dimensions{}
	}
	if s.btn.Clicked(gtx) {
		f := s.onAction
		s.msg = ""
		if f != nil {
			f()
		}
		return layout.Dimensions{}
	}
	gtx.Execute(op.InvalidateCmd{At: s.until})

	bg := theme.ColorSurfaceElevated
	if s.isError {
		bg = theme.ColorExpense
	}
	return theme.Card(gtx, bg, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(14), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th.Theme, s.msg)
					lbl.Color = theme.ColorTextPrimary
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if s.action == "" {
						return layout.Dimensions{}
					}
					btn := material.Button(th.Theme, &s.btn, s.action)
					btn.Background = bg
					btn.Color = theme.ColorPrimary
					btn.TextSize = unit.Sp(13)
					btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(10), Right: unit.Dp(10)}
					return btn.Layout(gtx)
				}),
			)
		})
	})
}
