package tracking

import "strings"

const TestCarrierCode = "test-carrier"
const TestCarrierLabel = "Test Carrier"

// IsTrackingMoreTestNumber reports whether the tracking number is a TrackingMore
// sandbox number (TEST1234… / TEST2234…). These only work with TestCarrierCode.
func IsTrackingMoreTestNumber(trackingNumber string) bool {
	tn := strings.ToUpper(strings.TrimSpace(trackingNumber))
	return strings.HasPrefix(tn, "TEST") && len(tn) >= 12
}
