package components

import (
	"image/color"
	"slices"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/shakilnwz/wallet/core/models"
	"github.com/shakilnwz/wallet/core/money"
	"github.com/shakilnwz/wallet/ui/theme"
)

var defaultExpenseCats = []string{
	"Food & Dining",
	"Housing & Rent",
	"Internet",
	"Utilities",
	"Groceries",
	"Transport",
	"Healthcare",
	"Shopping",
	"Personal",
	"Other",
}

var defaultIncomeCats = []string{
	"Salary",
	"Carry Forward",
	"Freelance",
	"Investments",
	"Bonus",
	"Gift",
	"Other",
}

// DefaultCategories returns the built-in categories for a transaction type.
func DefaultCategories(t models.TransactionType) []string {
	if t == models.TypeIncome {
		return defaultIncomeCats
	}
	return defaultExpenseCats
}

var scrimColor = color.NRGBA{A: 0xB0}

type QuickAddSheet struct {
	IsOpen bool

	// Categories returns the chips offered for a type, most relevant first.
	// When nil, DefaultCategories is used.
	Categories func(models.TransactionType) []string

	EditingID   string // set when editing an existing transaction
	Type        models.TransactionType
	AmountStr   string
	BenfEditor  widget.Editor
	Category    string
	Account     models.AccountType
	Date        time.Time
	NotesEditor widget.Editor
	ErrMsg      string

	cats []string

	// Controls
	TypeExpenseBtn widget.Clickable
	TypeIncomeBtn  widget.Clickable
	CloseBtn       widget.Clickable
	SubmitBtn      widget.Clickable

	AccountBankBtn widget.Clickable
	AccountCashBtn widget.Clickable

	TodayBtn     widget.Clickable
	YesterdayBtn widget.Clickable
	PrevDayBtn   widget.Clickable
	NextDayBtn   widget.Clickable

	CatBtns []widget.Clickable
	CatList layout.List
	Body    widget.List

	scrim widget.Clickable
	// sheet swallows taps on the sheet's background so they don't reach the
	// scrim underneath and close it.
	sheet widget.Clickable

	blurPending bool

	Keypad Keypad
}

func today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// Open shows an empty sheet for a new transaction.
func (s *QuickAddSheet) Open() {
	s.open()
	s.EditingID = ""
	s.AmountStr = ""
	s.Account = models.AccountBank
	s.Date = today()
	s.BenfEditor.SetText("")
	s.NotesEditor.SetText("")
	s.setType(models.TypeExpense, "")
}

// OpenEdit shows the sheet pre-filled with an existing transaction.
func (s *QuickAddSheet) OpenEdit(tx models.Transaction) {
	s.open()
	s.EditingID = tx.ID
	s.AmountStr = money.ToInput(tx.Amount)
	s.Account = tx.Account
	s.Date = today()
	if d, err := time.ParseInLocation("2006-01-02", tx.Date, time.Local); err == nil {
		s.Date = d
	}
	// Payee defaults to the category on save; don't show that back as typed text.
	benf := tx.Benf
	if benf == tx.ExpCat {
		benf = ""
	}
	s.BenfEditor.SetText(benf)
	s.NotesEditor.SetText(tx.Notes)
	s.setType(tx.Type, tx.ExpCat)
}

func (s *QuickAddSheet) open() {
	s.IsOpen = true
	s.ErrMsg = ""
	s.BenfEditor.SingleLine = true
	s.NotesEditor.SingleLine = true
	s.Body.Axis = layout.Vertical
	s.Body.Position = layout.Position{}
	s.CatList.Axis = layout.Horizontal
	s.CatList.Position = layout.Position{}
}

func (s *QuickAddSheet) Close() {
	s.IsOpen = false
	s.blurPending = true
}

func (s *QuickAddSheet) setType(t models.TransactionType, category string) {
	s.Type = t
	if s.Categories != nil {
		s.cats = s.Categories(t)
	} else {
		s.cats = DefaultCategories(t)
	}
	if category != "" && !slices.Contains(s.cats, category) {
		s.cats = append([]string{category}, s.cats...)
	}
	switch {
	case category != "":
		s.Category = category
	case len(s.cats) > 0:
		s.Category = s.cats[0]
	default:
		s.Category = "Other"
	}
	if len(s.CatBtns) < len(s.cats) {
		s.CatBtns = make([]widget.Clickable, len(s.cats))
	}
}

func (s *QuickAddSheet) Layout(
	gtx layout.Context,
	th *theme.Theme,
	onSubmit func(tx models.Transaction),
) layout.Dimensions {
	if s.blurPending {
		// Drop editor focus so Android hides the soft keyboard.
		s.blurPending = false
		gtx.Execute(key.FocusCmd{})
	}
	if !s.IsOpen {
		return layout.Dimensions{}
	}

	// Back (Android) and Escape (desktop) close the sheet. Only registered
	// while open, so Back still leaves the app otherwise.
	for {
		ev, ok := gtx.Event(key.Filter{Name: key.NameBack}, key.Filter{Name: key.NameEscape})
		if !ok {
			break
		}
		if e, ok := ev.(key.Event); ok && e.State == key.Press {
			s.Close()
		}
	}

	if s.CloseBtn.Clicked(gtx) || s.scrim.Clicked(gtx) {
		s.Close()
	}
	if !s.IsOpen {
		return layout.Dimensions{}
	}

	if s.TypeExpenseBtn.Clicked(gtx) && s.Type != models.TypeExpense {
		s.setType(models.TypeExpense, "")
	}
	if s.TypeIncomeBtn.Clicked(gtx) && s.Type != models.TypeIncome {
		s.setType(models.TypeIncome, "")
	}

	if s.AccountBankBtn.Clicked(gtx) {
		s.Account = models.AccountBank
	}
	if s.AccountCashBtn.Clicked(gtx) {
		s.Account = models.AccountCash
	}

	if s.TodayBtn.Clicked(gtx) {
		s.Date = today()
	}
	if s.YesterdayBtn.Clicked(gtx) {
		s.Date = today().AddDate(0, 0, -1)
	}
	if s.PrevDayBtn.Clicked(gtx) {
		s.Date = s.Date.AddDate(0, 0, -1)
	}
	if s.NextDayBtn.Clicked(gtx) && s.Date.Before(today()) {
		s.Date = s.Date.AddDate(0, 0, 1)
	}

	for i, c := range s.cats {
		if s.CatBtns[i].Clicked(gtx) {
			s.Category = c
		}
	}

	if s.SubmitBtn.Clicked(gtx) {
		amt, err := money.Parse(s.AmountStr)
		if err != nil || amt <= 0 {
			s.ErrMsg = "Enter an amount above ৳0"
		} else {
			benf := s.BenfEditor.Text()
			if benf == "" {
				benf = s.Category
			}
			onSubmit(models.Transaction{
				ID:      s.EditingID,
				Date:    s.Date.Format("2006-01-02"),
				Type:    s.Type,
				Benf:    benf,
				ExpCat:  s.Category,
				Amount:  amt,
				Account: s.Account,
				Notes:   s.NotesEditor.Text(),
			})
			s.Close()
			return layout.Dimensions{}
		}
	}

	// Dimmed backdrop over the whole screen; tapping it dismisses the sheet.
	gtx.Constraints.Min = gtx.Constraints.Max
	s.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		paint.FillShape(gtx.Ops, scrimColor, clip.Rect{Max: size}.Op())
		return layout.Dimensions{Size: size}
	})

	return layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		gtx.Constraints.Max.Y = gtx.Constraints.Max.Y * 92 / 100
		return s.sheet.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return theme.Card(gtx, theme.ColorSurface, unit.Dp(24), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return s.layoutBody(gtx, th)
				})
			})
		})
	})
}

func (s *QuickAddSheet) layoutBody(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	gap := func(dp unit.Dp) layout.Widget { return layout.Spacer{Height: dp}.Layout }

	rows := []layout.Widget{
		// Header: Title & Close Button
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txt := "Add Transaction"
					if s.EditingID != "" {
						txt = "Edit Transaction"
					}
					title := material.H6(th.Theme, txt)
					title.Color = theme.ColorTextPrimary
					title.TextSize = unit.Sp(16)
					return title.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.CloseBtn, "X", false, theme.ColorSurfaceElevated, unit.Sp(13))
				}),
			)
		},
		gap(12),

		// Type Segmented Switch: Expense vs Income
		func(gtx layout.Context) layout.Dimensions {
			return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return segment(gtx, th, &s.TypeExpenseBtn, "Expense", s.Type == models.TypeExpense, theme.ColorExpense)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return segment(gtx, th, &s.TypeIncomeBtn, "Income", s.Type == models.TypeIncome, theme.ColorIncome)
						}),
					)
				})
			})
		},
		gap(10),

		// Live Amount Display
		func(gtx layout.Context) layout.Dimensions {
			return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					valText := s.AmountStr
					if valText == "" {
						valText = "0"
					}
					lbl := material.H5(th.Theme, "৳ "+valText)
					lbl.Color = theme.ColorExpense
					if s.Type == models.TypeIncome {
						lbl.Color = theme.ColorIncome
					}
					lbl.TextSize = unit.Sp(26)
					return lbl.Layout(gtx)
				})
			})
		},
		gap(10),

		// Category chips: horizontally scrollable so none are cut off.
		func(gtx layout.Context) layout.Dimensions {
			return s.CatList.Layout(gtx, len(s.cats), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.Inset{Right: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.CatBtns[i], s.cats[i], s.Category == s.cats[i], theme.ColorPrimary, unit.Sp(12))
				})
			})
		},
		gap(10),

		// Payee & notes
		func(gtx layout.Context) layout.Dimensions {
			hint := "Paid to (optional)"
			if s.Type == models.TypeIncome {
				hint = "Received from (optional)"
			}
			return textField(gtx, th, &s.BenfEditor, hint)
		},
		gap(8),
		func(gtx layout.Context) layout.Dimensions {
			return textField(gtx, th, &s.NotesEditor, "Note (optional)")
		},
		gap(10),

		// Date: step by day, with Today / Yesterday shortcuts.
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.PrevDayBtn, "<", false, theme.ColorSurfaceElevated, unit.Sp(13))
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th.Theme, theme.FormatDateDisplay(s.Date.Format("2006-01-02"))+"  ·  "+s.Date.Format("Mon"))
						lbl.Color = theme.ColorTextPrimary
						return lbl.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.NextDayBtn, ">", false, theme.ColorSurfaceElevated, unit.Sp(13))
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.TodayBtn, "Today", s.Date.Equal(today()), theme.ColorPrimary, unit.Sp(12))
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return pill(gtx, th, &s.YesterdayBtn, "Yesterday", s.Date.Equal(today().AddDate(0, 0, -1)), theme.ColorPrimary, unit.Sp(12))
				}),
			)
		},
		gap(10),

		// Account Toggle: Bank Account vs Cash
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return segment(gtx, th, &s.AccountBankBtn, "Account (Bank/MFS)", s.Account == models.AccountBank, theme.ColorPrimary)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return segment(gtx, th, &s.AccountCashBtn, "Cash", s.Account == models.AccountCash, theme.ColorIncome)
				}),
			)
		},
		gap(10),

		// Embedded Numeric Keypad, kept near the bottom for one-thumb entry.
		func(gtx layout.Context) layout.Dimensions {
			return s.Keypad.Layout(gtx, th, s.AmountStr, func(newVal string) {
				s.AmountStr = newVal
				s.ErrMsg = ""
			})
		},
		gap(10),

		func(gtx layout.Context) layout.Dimensions {
			if s.ErrMsg == "" {
				return layout.Dimensions{}
			}
			lbl := material.Caption(th.Theme, s.ErrMsg)
			lbl.Color = theme.ColorExpense
			lbl.TextSize = unit.Sp(12)
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, lbl.Layout)
		},

		// Submit Action Button
		func(gtx layout.Context) layout.Dimensions {
			submitText := "Save Expense"
			submitBg := theme.ColorExpense
			if s.Type == models.TypeIncome {
				submitText = "Save Income"
				submitBg = theme.ColorIncome
			}
			if s.EditingID != "" {
				submitText = "Save Changes"
			}
			btn := material.Button(th.Theme, &s.SubmitBtn, submitText)
			btn.Background = submitBg
			btn.Color = theme.ColorTextPrimary
			btn.CornerRadius = unit.Dp(12)
			btn.TextSize = unit.Sp(15)
			btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
			return btn.Layout(gtx)
		},
	}

	return material.List(th.Theme, &s.Body).Layout(gtx, len(rows), func(gtx layout.Context, i int) layout.Dimensions {
		return rows[i](gtx)
	})
}

// pill is a small rounded button that highlights with activeBg when active.
func pill(gtx layout.Context, th *theme.Theme, click *widget.Clickable, label string, active bool, activeBg color.NRGBA, size unit.Sp) layout.Dimensions {
	btn := material.Button(th.Theme, click, label)
	btn.Background = theme.ColorSurfaceElevated
	btn.Color = theme.ColorTextMuted
	if active {
		btn.Background = activeBg
		btn.Color = theme.ColorTextPrimary
	}
	btn.CornerRadius = unit.Dp(8)
	btn.TextSize = size
	btn.Inset = layout.Inset{Top: unit.Dp(7), Bottom: unit.Dp(7), Left: unit.Dp(10), Right: unit.Dp(10)}
	return btn.Layout(gtx)
}

// segment is a full-width toggle button used in segmented controls.
func segment(gtx layout.Context, th *theme.Theme, click *widget.Clickable, label string, active bool, activeBg color.NRGBA) layout.Dimensions {
	btn := material.Button(th.Theme, click, label)
	btn.Background = theme.ColorSurfaceElevated
	btn.Color = theme.ColorTextMuted
	if active {
		btn.Background = activeBg
		btn.Color = theme.ColorTextPrimary
	}
	btn.CornerRadius = unit.Dp(10)
	btn.TextSize = unit.Sp(13)
	btn.Inset = layout.Inset{Top: unit.Dp(9), Bottom: unit.Dp(9)}
	return btn.Layout(gtx)
}

func textField(gtx layout.Context, th *theme.Theme, ed *widget.Editor, hint string) layout.Dimensions {
	return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			e := material.Editor(th.Theme, ed, hint)
			e.Color = theme.ColorTextPrimary
			e.HintColor = theme.ColorTextDim
			e.TextSize = unit.Sp(14)
			return e.Layout(gtx)
		})
	})
}
