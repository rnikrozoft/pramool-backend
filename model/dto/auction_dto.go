package dto

import "time"

type CreateAuctionRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
	StartPrice  int64  `json:"start_price"`
	BidStep     int64  `json:"bid_step"`
	EndAt       string `json:"end_at"`
}

type CreateAuctionResponse struct {
	AuctionID string `json:"auction_id"`
}

type SellerAuctionItem struct {
	AuctionID     string `json:"auction_id"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	Status        string `json:"status"`
	StartPrice    int64  `json:"start_price"`
	CurrentBid    int64  `json:"current_bid"`
	TotalBids     int64  `json:"total_bids"`
	EndAt         string `json:"end_at"`
	CoverImageURL string `json:"cover_image_url"`
}

// SellerEarningItem is one row in GET /seller/earnings.
type SellerEarningItem struct {
	EarningID    int64     `json:"earning_id"`
	AuctionID    string    `json:"auction_id"`
	WinnerUserID string    `json:"winner_user_id"`
	Amount       int64     `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
