package timing

import "testing"

func TestSupportedFixedRates(t *testing.T) {
	for input, want := range map[int]int{0: 60, 50: 50, 60: 60} {
		got, err := Normalize(input)
		if err != nil || got != want {
			t.Fatalf("Normalize(%d) = %d, %v", input, got, err)
		}
	}
	for _, input := range []int{-1, 30, 120, 144} {
		if _, err := Normalize(input); err == nil {
			t.Fatalf("accepted unsupported rate %d", input)
		}
	}
}
