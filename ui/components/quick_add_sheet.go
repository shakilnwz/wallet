package components

import (
	"fmt"
	"time"

	"gioui.org/layout"
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

type QuickAddSheet struct {
	IsOpen bool

	Type        models.TransactionType // income | expense
	AmountStr   string
	BenfEditor  widget.Editor
	Category    string
	Account     models.AccountType
	DateStr     string // YYYY-MM-DD
	NotesEditor widget.Editor

	// Controls
	TypeExpenseBtn widget.Clickable
	TypeIncomeBtn  widget.Clickable
	CloseBtn       widget.Clickable
	SubmitBtn      widget.Clickable

	AccountBankBtn widget.Clickable
	AccountCashBtn widget.Clickable

	TodayBtn     widget.Clickable
	YesterdayBtn widget.Clickable

	CatBtns [10]widget.Clickable

	Keypad Keypad
}

func (s *QuickAddSheet) Open() {
	s.IsOpen = true
	s.Type = models.TypeExpense
	s.AmountStr = ""
	s.Category = "Food & Dining"
	s.Account = models.AccountBank
	s.DateStr = time.Now().Format("2006-01-02")
	s.BenfEditor.SetText("")
	s.NotesEditor.SetText("")
}

func (s *QuickAddSheet) Close() {
	s.IsOpen = false
}

func (s *QuickAddSheet) Layout(
	gtx layout.Context,
	th *theme.Theme,
	onSubmit func(tx models.Transaction),
) layout.Dimensions {
	if !s.IsOpen {
		return layout.Dimensions{}
	}

	if s.CloseBtn.Clicked(gtx) {
		s.Close()
		return layout.Dimensions{}
	}

	if s.TypeExpenseBtn.Clicked(gtx) {
		s.Type = models.TypeExpense
		s.Category = "Food & Dining"
	}
	if s.TypeIncomeBtn.Clicked(gtx) {
		s.Type = models.TypeIncome
		s.Category = "Salary"
	}

	if s.AccountBankBtn.Clicked(gtx) {
		s.Account = models.AccountBank
	}
	if s.AccountCashBtn.Clicked(gtx) {
		s.Account = models.AccountCash
	}

	if s.TodayBtn.Clicked(gtx) {
		s.DateStr = time.Now().Format("2006-01-02")
	}
	if s.YesterdayBtn.Clicked(gtx) {
		s.DateStr = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	// Categories selection
	cats := defaultExpenseCats
	if s.Type == models.TypeIncome {
		cats = defaultIncomeCats
	}
	for i, c := range cats {
		if i < len(s.CatBtns) && s.CatBtns[i].Clicked(gtx) {
			s.Category = c
		}
	}

	// Submit handling
	if s.SubmitBtn.Clicked(gtx) {
		amt, err := money.Parse(s.AmountStr)
		if err == nil && amt > 0 {
			benf := s.BenfEditor.Text()
			if benf == "" {
				benf = s.Category
			}
			onSubmit(models.Transaction{
				Date:    s.DateStr,
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

	// Modal Overlay Layout
	return theme.Card(gtx, theme.ColorSurface, unit.Dp(24), func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    unit.Dp(16),
			Bottom: unit.Dp(16),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				// Header: Title & Close Button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H6(th.Theme, "Quick Add Transaction")
							title.Color = theme.ColorTextPrimary
							title.TextSize = unit.Sp(16)
							return title.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th.Theme, &s.CloseBtn, "X")
							btn.Background = theme.ColorSurfaceElevated
							btn.Color = theme.ColorTextMuted
							btn.CornerRadius = unit.Dp(16)
							btn.TextSize = unit.Sp(13)
							btn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}
							return btn.Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),

				// Type Segmented Switch: Expense vs Income
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									bg := theme.ColorSurfaceElevated
									fg := theme.ColorTextMuted
									if s.Type == models.TypeExpense {
										bg = theme.ColorExpense
										fg = theme.ColorTextPrimary
									}
									btn := material.Button(th.Theme, &s.TypeExpenseBtn, "Expense")
									btn.Background = bg
									btn.Color = fg
									btn.CornerRadius = unit.Dp(8)
									btn.TextSize = unit.Sp(13)
									btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}
									return btn.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									bg := theme.ColorSurfaceElevated
									fg := theme.ColorTextMuted
									if s.Type == models.TypeIncome {
										bg = theme.ColorIncome
										fg = theme.ColorTextPrimary
									}
									btn := material.Button(th.Theme, &s.TypeIncomeBtn, "Income")
									btn.Background = bg
									btn.Color = fg
									btn.CornerRadius = unit.Dp(8)
									btn.TextSize = unit.Sp(13)
									btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}
									return btn.Layout(gtx)
								}),
							)
						})
					})
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Live Amount Display
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return theme.Card(gtx, theme.ColorSurfaceElevated, unit.Dp(14), func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							valText := s.AmountStr
							if valText == "" {
								valText = "0.00"
							}
							display := fmt.Sprintf("৳ %s", valText)
							lbl := material.H5(th.Theme, display)
							if s.Type == models.TypeIncome {
								lbl.Color = theme.ColorIncome
							} else {
								lbl.Color = theme.ColorExpense
							}
							lbl.TextSize = unit.Sp(24)
							return lbl.Layout(gtx)
						})
					})
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Category selection chips row
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							var chipChildren []layout.FlexChild
							for i, c := range cats {
								if i >= len(s.CatBtns) {
									break
								}
								catName := c
								btnIdx := i
								chipChildren = append(chipChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									bg := theme.ColorSurfaceElevated
									fg := theme.ColorTextMuted
									if s.Category == catName {
										bg = theme.ColorPrimary
										fg = theme.ColorTextPrimary
									}
									btn := material.Button(th.Theme, &s.CatBtns[btnIdx], catName)
									btn.Background = bg
									btn.Color = fg
									btn.CornerRadius = unit.Dp(8)
									btn.TextSize = unit.Sp(11)
									btn.Inset = layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}
									return layout.Inset{Right: unit.Dp(6)}.Layout(gtx, btn.Layout)
								}))
							}
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, chipChildren...)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Embedded Numeric Keypad
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.Keypad.Layout(gtx, th, s.AmountStr, func(newVal string) {
						s.AmountStr = newVal
					})
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Account Toggle: Bank Account vs Cash
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							bg := theme.ColorSurfaceElevated
							fg := theme.ColorTextMuted
							if s.Account == models.AccountBank {
								bg = theme.ColorPrimary
								fg = theme.ColorTextPrimary
							}
							btn := material.Button(th.Theme, &s.AccountBankBtn, "Account (Bank/MFS)")
							btn.Background = bg
							btn.Color = fg
							btn.CornerRadius = unit.Dp(10)
							btn.TextSize = unit.Sp(12)
							btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}
							return btn.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							bg := theme.ColorSurfaceElevated
							fg := theme.ColorTextMuted
							if s.Account == models.AccountCash {
								bg = theme.ColorIncome
								fg = theme.ColorTextPrimary
							}
							btn := material.Button(th.Theme, &s.AccountCashBtn, "Cash")
							btn.Background = bg
							btn.Color = fg
							btn.CornerRadius = unit.Dp(10)
							btn.TextSize = unit.Sp(12)
							btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}
							return btn.Layout(gtx)
						}),
					)
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),

				// Submit Action Button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					submitText := "Save Expense"
					submitBg := theme.ColorExpense
					if s.Type == models.TypeIncome {
						submitText = "Save Income"
						submitBg = theme.ColorIncome
					}
					btn := material.Button(th.Theme, &s.SubmitBtn, submitText)
					btn.Background = submitBg
					btn.Color = theme.ColorTextPrimary
					btn.CornerRadius = unit.Dp(12)
					btn.TextSize = unit.Sp(15)
					btn.Inset = layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}
					return btn.Layout(gtx)
				}),
			)
		})
	})
}
