package entity

import "time"

// SellerEarning is a seller_earnings row.
type SellerEarning struct {
	EarningID    int64     `db:"earning_id"`
	AuctionID    string    `db:"auction_id"`
	WinnerUserID string    `db:"winner_user_id"`
	Amount       int64     `db:"amount"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
}
