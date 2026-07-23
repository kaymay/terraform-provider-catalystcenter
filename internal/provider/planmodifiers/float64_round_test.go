package planmodifiers

import "testing"

func TestRoundTo5(t *testing.T) {
	cases := map[float64]float64{
		51.730618:  51.73062,
		15.826622:  15.82662,
		44.225605:  44.22561,
		6.140181:   6.14018,
		-121.83201: -121.83201,
		37.338:     37.338,
	}
	for in, want := range cases {
		if got := RoundTo(in, 5); got != want {
			t.Errorf("RoundTo(%v,5)=%v want %v", in, got, want)
		}
	}
}
