package money

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"0", 0, false},
		{"1500", 150000, false},
		{"12.5", 1250, false},
		{"12.05", 1205, false},
		{"0.1", 10, false},
		{".5", 50, false},
		{"7.", 700, false},
		{" 42 ", 4200, false},
		{"", 0, true},
		{"abc", 0, true},
		{"1.234", 0, true},
		{"-5", 0, true},
		{"1.-5", 0, true},
		{"1.2.3", 0, true},
		{".", 0, true},
		{"+5", 0, true},
		{"99999999999999999999", 0, true},
	}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("Parse(%q) = %d, %v; want %d, err=%v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestToInputRoundTrip(t *testing.T) {
	for _, p := range []int64{0, 5, 10, 100, 1250, 150000} {
		got, err := Parse(ToInput(p))
		if err != nil || got != p {
			t.Errorf("round trip %d -> %q -> %d (%v)", p, ToInput(p), got, err)
		}
	}
}
