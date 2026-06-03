package tracking

import "testing"

func TestIsTrackingMoreTestNumber(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"TEST1234123421", true},
		{"TEST2234123441", true},
		{"test1234123401", true},
		{" TH01283G4AQV8A ", false},
		{"TEST123", false},
		{"", false},
	}
	for _, tc := range tests {
		if got := IsTrackingMoreTestNumber(tc.in); got != tc.want {
			t.Errorf("IsTrackingMoreTestNumber(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
