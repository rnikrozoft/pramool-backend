package repository

import (
	"context"
	"time"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

// LockAuctionBySellerForUpdate loads the auction row for the seller with FOR UPDATE.
func (r auctionRepo) LockAuctionBySellerForUpdate(ctx context.Context, tx bun.Tx, auctionID, sellerID string) (*entity.Auction, error) {
	item := new(entity.Auction)
	err := tx.NewRaw(`
		SELECT auction_id, seller_id, title, category, item_condition AS condition, description,
			start_price, bid_step, current_bid, total_bids, status, end_at,
			COALESCE(allow_early_close, FALSE) AS allow_early_close,
			COALESCE(early_close_hold_amount, 0) AS early_close_hold_amount,
			cover_image_url,
			created_at, updated_at, COALESCE(winner_id, '') AS winner_id
		FROM auctions
		WHERE auction_id = ? AND seller_id = ?
		FOR UPDATE
	`, auctionID, sellerID).Scan(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r auctionRepo) CountAuctionBidsTx(ctx context.Context, tx bun.Tx, auctionID string) (int64, error) {
	var n int64
	err := tx.NewRaw(`SELECT COUNT(*)::bigint FROM auction_bids WHERE auction_id = ?`, auctionID).Scan(ctx, &n)
	return n, err
}

func (r auctionRepo) CountHeldBidHoldsTx(ctx context.Context, tx bun.Tx, auctionID string) (int64, error) {
	var n int64
	err := tx.NewRaw(`SELECT COUNT(*)::bigint FROM auction_bid_holds WHERE auction_id = ? AND hold_status = 'held'`, auctionID).Scan(ctx, &n)
	return n, err
}

// ApplyAuctionReopenTx sets the auction back to active; WHERE clause must match reopen preconditions.
func (r auctionRepo) ApplyAuctionReopenTx(ctx context.Context, tx bun.Tx, auctionID, sellerID string, endAt time.Time) (int64, error) {
	res, err := tx.NewRaw(`
		UPDATE auctions SET
			status = 'active',
			current_bid = start_price,
			total_bids = 0,
			end_at = ?,
			winner_id = NULL,
			settled_at = NULL,
			early_close_hold_amount = 0,
			updated_at = NOW()
		WHERE auction_id = ? AND seller_id = ?
		  AND status = 'closed'
		  AND total_bids = 0
		  AND (winner_id IS NULL OR winner_id = '')
	`, endAt, auctionID, sellerID).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
