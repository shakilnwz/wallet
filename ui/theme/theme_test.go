package theme

import "testing"

func TestFormatCurrency(t *testing.T) {
	tests := map[int64]string{
		0:          "৳0.00",
		5:          "৳0.05",
		30:         "৳0.30",
		123456:     "৳1,234.56",
		1234550:    "৳12,345.50",
		1000000000: "৳1,00,00,000.00",
		-150050:    "-৳1,500.50",
	}
	for paisa, want := range tests {
		if got := FormatCurrency(paisa); got != want {
			t.Errorf("FormatCurrency(%d) = %q, want %q", paisa, got, want)
		}
	}
}
