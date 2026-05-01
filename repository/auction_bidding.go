package repository

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
)

func (r auctionRepo) ListSellerEarnings(ctx context.Context, sellerID string, limit, offset int) ([]entity.SellerEarning, error) {
	rows, err := r.bun.QueryContext(ctx, `
		SELECT earning_id, auction_id, winner_user_id, amount, status, created_at
		FROM seller_earnings
		WHERE seller_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, sellerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.SellerEarning, 0)
	for rows.Next() {
		var item entity.SellerEarning
		if err := rows.Scan(&item.EarningID, &item.AuctionID, &item.WinnerUserID, &item.Amount, &item.Status, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
