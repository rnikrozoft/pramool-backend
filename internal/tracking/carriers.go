package tracking

import "strings"

// Carrier is a supported Thailand courier for TrackingMore API v2.
type Carrier struct {
	Code        string `json:"code"`
	Label       string `json:"label"`
	TrackURLFmt string `json:"track_url_fmt,omitempty"`
}

var thailandCarriers = []Carrier{
	{Code: "flashexpress", Label: "Flash Express", TrackURLFmt: "https://www.flashexpress.co.th/tracking/?se=%s"},
	{Code: "kerryexpress-th", Label: "Kerry Express", TrackURLFmt: "https://th.kerryexpress.com/en/track/?track=%s"},
	{Code: "jt-express-th", Label: "J&T Express", TrackURLFmt: "https://www.jtexpress.co.th/service/track?waybillNo=%s"},
	{Code: "thailand-post", Label: "ไปรษณีย์ไทย", TrackURLFmt: "https://track.thailandpost.co.th/?trackNumber=%s"},
}

func ListCarriers() []Carrier {
	out := make([]Carrier, len(thailandCarriers))
	copy(out, thailandCarriers)
	return out
}

func CarrierByCode(code string) (Carrier, bool) {
	code = strings.ToLower(strings.TrimSpace(code))
	for _, c := range thailandCarriers {
		if c.Code == code {
			return c, true
		}
	}
	return Carrier{}, false
}

func CarrierLabel(code string) string {
	if c, ok := CarrierByCode(code); ok {
		return c.Label
	}
	return strings.TrimSpace(code)
}

func KnownCarrierCodes() map[string]struct{} {
	out := make(map[string]struct{}, len(thailandCarriers))
	for _, c := range thailandCarriers {
		out[c.Code] = struct{}{}
	}
	return out
}
