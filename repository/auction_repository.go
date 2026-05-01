package repository

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

type AuctionRepository interface {
	BeginTx(ctx context.Context) (bun.Tx, error)
	CreateAuctionWithTx(ctx context.Context, tx bun.Tx, auction entity.Auction) error
	CreateAuctionImagesWithTx(ctx context.Context, tx bun.Tx, images []entity.AuctionImage) error
	ListAuctionsBySellerID(ctx context.Context, sellerID string) ([]entity.Auction, error)
	ListSellerEarnings(ctx context.Context, sellerID string, limit, offset int) ([]entity.SellerEarning, error)
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
		start_price, bid_step, current_bid, total_bids, status, end_at, cover_image_url
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		auction.CoverImageURL,
	).Exec(ctx)
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
		start_price, bid_step, current_bid, total_bids, status, end_at, cover_image_url,
		created_at, updated_at
	FROM auctions
	WHERE seller_id = ?
	ORDER BY created_at DESC
	`
	err := r.bun.NewRaw(query, sellerID).Scan(ctx, &items)
	return items, err
}
