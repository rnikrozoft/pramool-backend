package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

type AuctionService interface {
	CreateAuction(ctx context.Context, sellerID string, req dto.CreateAuctionRequest, imagePaths []string) (*dto.CreateAuctionResponse, error)
	GetSellerAuctions(ctx context.Context, sellerID string) ([]dto.SellerAuctionItem, error)
	ListSellerEarnings(ctx context.Context, sellerID string, limit, offset int) ([]dto.SellerEarningItem, error)
}

type auction struct {
	repo repository.AuctionRepository
}

func NewAuctionService(repo repository.AuctionRepository) AuctionService {
	return auction{repo: repo}
}

func (s auction) CreateAuction(ctx context.Context, sellerID string, req dto.CreateAuctionRequest, imagePaths []string) (*dto.CreateAuctionResponse, error) {
	if strings.TrimSpace(sellerID) == "" {
		return nil, fmt.Errorf("missing seller id")
	}
	if strings.TrimSpace(req.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.StartPrice <= 0 || req.BidStep <= 0 {
		return nil, fmt.Errorf("invalid price settings")
	}
	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}

	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return nil, fmt.Errorf("invalid end_at format")
	}
	auctionID := generateAuctionID()

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	mainAuction := entity.Auction{
		AuctionID:     auctionID,
		SellerID:      sellerID,
		Title:         strings.TrimSpace(req.Title),
		Category:      strings.TrimSpace(req.Category),
		Condition:     strings.TrimSpace(req.Condition),
		Description:   strings.TrimSpace(req.Description),
		StartPrice:    req.StartPrice,
		BidStep:       req.BidStep,
		CurrentBid:    req.StartPrice,
		TotalBids:     0,
		Status:        "active",
		EndAt:         endAt,
		CoverImageURL: imagePaths[0],
	}
	if err := s.repo.CreateAuctionWithTx(ctx, tx, mainAuction); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	images := make([]entity.AuctionImage, 0, len(imagePaths))
	for i, p := range imagePaths {
		images = append(images, entity.AuctionImage{
			AuctionID: auctionID,
			ImageURL:  p,
			SortOrder: i,
		})
	}
	if err := s.repo.CreateAuctionImagesWithTx(ctx, tx, images); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &dto.CreateAuctionResponse{AuctionID: auctionID}, nil
}

func (s auction) GetSellerAuctions(ctx context.Context, sellerID string) ([]dto.SellerAuctionItem, error) {
	items, err := s.repo.ListAuctionsBySellerID(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SellerAuctionItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.SellerAuctionItem{
			AuctionID:     item.AuctionID,
			Title:         item.Title,
			Category:      item.Category,
			Status:        item.Status,
			StartPrice:    item.StartPrice,
			CurrentBid:    item.CurrentBid,
			TotalBids:     item.TotalBids,
			EndAt:         item.EndAt.Format(time.RFC3339),
			CoverImageURL: item.CoverImageURL,
		})
	}
	return result, nil
}

func (s auction) ListSellerEarnings(ctx context.Context, sellerID string, limit, offset int) ([]dto.SellerEarningItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.repo.ListSellerEarnings(ctx, sellerID, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]dto.SellerEarningItem, 0, len(items))
	for _, e := range items {
		out = append(out, dto.SellerEarningItem{
			EarningID:    e.EarningID,
			AuctionID:    e.AuctionID,
			WinnerUserID: e.WinnerUserID,
			Amount:       e.Amount,
			Status:       e.Status,
			CreatedAt:    e.CreatedAt,
		})
	}
	return out, nil
}

func generateAuctionID() string {
	return fmt.Sprintf("AUC-%s", strings.ToUpper(strings.ReplaceAll(time.Now().Format("20060102-150405.000"), ".", "")))
}
