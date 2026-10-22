// Package money handles BDT amounts stored as int64 paisa (1 ৳ = 100 paisa).
package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalid = errors.New("invalid amount")

// Parse converts user input such as "1500", "12.5" or "0.05" into paisa.
// Negative values and more than two decimal places are rejected.
func Parse(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, ErrInvalid
	}

	whole, frac, _ := strings.Cut(s, ".")
	if whole == "" {
		whole = "0"
	}
	if len(frac) > 2 || s == "." {
		return 0, ErrInvalid
	}

	w, err := strconv.ParseUint(whole, 10, 63)
	if err != nil {
		return 0, ErrInvalid
	}

	var f uint64
	if frac != "" {
		f, err = strconv.ParseUint(frac, 10, 8)
		if err != nil {
			return 0, ErrInvalid
		}
		if len(frac) == 1 {
			f *= 10
		}
	}

	if w > (1<<63-1-f)/100 {
		return 0, ErrInvalid
	}
	return int64(w*100 + f), nil
}

// ToInput renders paisa back into keypad text: "1500", "12.50".
func ToInput(paisa int64) string {
	if paisa%100 == 0 {
		return strconv.FormatInt(paisa/100, 10)
	}
	return fmt.Sprintf("%d.%02d", paisa/100, paisa%100)
}
