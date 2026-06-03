package dto

type MarkShippedRequest struct {
	CarrierCode    string `json:"carrier_code,omitempty"`
	TrackingNumber string `json:"tracking_number"`
}

type ShipmentCarrierItem struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type ShipmentTrackingEvent struct {
	Status     string `json:"status"`
	Location   string `json:"location,omitempty"`
	Note       string `json:"note,omitempty"`
	OccurredAt string `json:"occurred_at,omitempty"`
}

type ShipmentTrackingResponse struct {
	AuctionID      string                  `json:"auction_id"`
	CarrierCode    string                  `json:"carrier_code"`
	CarrierName    string                  `json:"carrier_name"`
	TrackingNumber string                  `json:"tracking_number"`
	ShipmentStatus string                  `json:"shipment_status"`
	TrackURL       string                  `json:"track_url,omitempty"`
	Events         []ShipmentTrackingEvent `json:"events"`
	CanConfirm     bool                    `json:"can_confirm"`
}

// WinnerShippingAddressResponse is buyer contact info for seller before mark-shipped.
type WinnerShippingAddressResponse struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Phone             string `json:"phone"`
	Address           string `json:"address"`
	DisclosureNotice  string `json:"disclosure_notice"`
	PhoneFullyVisible bool   `json:"phone_fully_visible"`
}
