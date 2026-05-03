package repository

import (
	"context"
	"time"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

type AuctionRepository interface {
	BeginTx(ctx context.Context) (bun.Tx, error)
	CreateAuctionWithTx(ctx context.Context, tx bun.Tx, auction entity.Auction) error
	CreateAuctionImagesWithTx(ctx context.Context, tx bun.Tx, images []entity.AuctionImage) error
	ListAuctionsBySellerID(ctx context.Context, sellerID string) ([]entity.Auction, error)

	LockAuctionBySellerForUpdate(ctx context.Context, tx bun.Tx, auctionID, sellerID string) (*entity.Auction, error)
	CountAuctionBidsTx(ctx context.Context, tx bun.Tx, auctionID string) (int64, error)
	CountHeldBidHoldsTx(ctx context.Context, tx bun.Tx, auctionID string) (int64, error)
	ApplyAuctionReopenTx(ctx context.Context, tx bun.Tx, auctionID, sellerID string, endAt time.Time) (int64, error)
	InsertListingDepositHoldTx(ctx context.Context, tx bun.Tx, sellerID, auctionID string, holdAmount, balanceBefore, balanceAfter int64, note string) error
}

type auctionRepo struct {
	bun *bun.DB
}

func NewAuctionRepository(bun *bun.DB) AuctionRepository {
	return auctionRepo{bun: bun}
}

func (r auctionRepo) BeginTx(ctx context.Context) (bun.Tx, error) {
	return r.bun.BeginTx(ctx, nil)
}

func (r auctionRepo) CreateAuctionWithTx(ctx context.Context, tx bun.Tx, auction entity.Auction) error {
	query := `
	INSERT INTO auctions (
		auction_id, seller_id, title, category, item_condition, description,
		start_price, bid_step, current_bid, total_bids, status, end_at, allow_early_close, early_close_hold_amount, cover_image_url
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.NewRaw(query,
		auction.AuctionID,
		auction.SellerID,
		auction.Title,
		auction.Category,
		auction.Condition,
		auction.Description,
		auction.StartPrice,
		auction.BidStep,
		auction.CurrentBid,
		auction.TotalBids,
		auction.Status,
		auction.EndAt,
		auction.AllowEarlyClose,
		auction.EarlyCloseHoldAmount,
		auction.CoverImageURL,
	).Exec(ctx)
	return err
}

func (r auctionRepo) InsertListingDepositHoldTx(ctx context.Context, tx bun.Tx, sellerID, auctionID string, holdAmount, balanceBefore, balanceAfter int64, note string) error {
	if holdAmount <= 0 {
		return nil
	}
	ledgerDelta := -holdAmount
	query := `
		INSERT INTO bid_transactions (user_id, auction_id, tx_type, amount, balance_before, balance_after, note, bid_amount)
		VALUES (?, ?, 'listing_deposit_hold', ?, ?, ?, ?, ?)
	`
	_, err := tx.NewRaw(query, sellerID, auctionID, ledgerDelta, balanceBefore, balanceAfter, note, holdAmount).Exec(ctx)
	return err
}

func (r auctionRepo) CreateAuctionImagesWithTx(ctx context.Context, tx bun.Tx, images []entity.AuctionImage) error {
	query := `INSERT INTO auction_images (auction_id, image_url, sort_order) VALUES (?, ?, ?)`
	for _, img := range images {
		if _, err := tx.NewRaw(query, img.AuctionID, img.ImageURL, img.SortOrder).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r auctionRepo) ListAuctionsBySellerID(ctx context.Context, sellerID string) ([]entity.Auction, error) {
	items := make([]entity.Auction, 0)
	query := `
	SELECT auction_id, seller_id, title, category, item_condition AS condition, description,
		start_price, bid_step, current_bid, total_bids, status, end_at, COALESCE(allow_early_close, FALSE) AS allow_early_close, cover_image_url,
		created_at, updated_at
	FROM auctions
	WHERE seller_id = ?
	ORDER BY created_at DESC
	`
	err := r.bun.NewRaw(query, sellerID).Scan(ctx, &items)
	return items, err
}
