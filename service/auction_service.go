package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

// ErrAuctionReopenNotAllowed is returned when reopen preconditions are not met.
var ErrAuctionReopenNotAllowed = errors.New("auction cannot be reopened: must be closed with no bids")

type AuctionService interface {
	CreateAuction(ctx context.Context, sellerID string, req dto.CreateAuctionRequest, imagePaths []string) (*dto.CreateAuctionResponse, error)
	GetSellerAuctions(ctx context.Context, sellerID string) ([]dto.SellerAuctionItem, error)
	ReopenAuctionNoBids(ctx context.Context, sellerID, auctionID, endAtRFC3339 string) error
}

type auction struct {
	repo     repository.AuctionRepository
	userRepo repository.UserRepository
}

func NewAuctionService(repo repository.AuctionRepository, userRepo repository.UserRepository) AuctionService {
	return auction{repo: repo, userRepo: userRepo}
}

func (s auction) CreateAuction(ctx context.Context, sellerID string, req dto.CreateAuctionRequest, imagePaths []string) (*dto.CreateAuctionResponse, error) {
	if strings.TrimSpace(sellerID) == "" {
		return nil, fmt.Errorf("missing seller id")
	}
	if strings.TrimSpace(req.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.StartPrice < 100 || req.BidStep <= 0 {
		return nil, fmt.Errorf("invalid price settings")
	}
	if len(imagePaths) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}

	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return nil, fmt.Errorf("invalid end_at format")
	}
	if !endAt.After(time.Now()) {
		return nil, fmt.Errorf("end_at must be in the future")
	}
	auctionID := generateAuctionID()

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	ok, balBefore, balAfter, err := s.userRepo.DeductListingDepositTx(ctx, tx, sellerID, req.StartPrice)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if !ok {
		_ = tx.Rollback()
		return nil, fmt.Errorf("insufficient credit for start price (%d THB required)", req.StartPrice)
	}

	mainAuction := entity.Auction{
		AuctionID:            auctionID,
		SellerID:             sellerID,
		Title:                strings.TrimSpace(req.Title),
		Category:             strings.TrimSpace(req.Category),
		Condition:            strings.TrimSpace(req.Condition),
		Description:          strings.TrimSpace(req.Description),
		StartPrice:           req.StartPrice,
		BidStep:              req.BidStep,
		CurrentBid:           req.StartPrice,
		TotalBids:            0,
		Status:               "active",
		EndAt:                endAt,
		AllowEarlyClose:      req.AllowEarlyClose,
		EarlyCloseHoldAmount: 0,
		CoverImageURL:        imagePaths[0],
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

	if err := s.repo.InsertListingDepositHoldTx(ctx, tx, sellerID, auctionID, req.StartPrice, balBefore, balAfter, "หักมัดจำประกาศเมื่อสร้างประมูล"); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &dto.CreateAuctionResponse{AuctionID: auctionID}, nil
}

func (s auction) ReopenAuctionNoBids(ctx context.Context, sellerID, auctionID, endAtRFC3339 string) error {
	if strings.TrimSpace(sellerID) == "" || strings.TrimSpace(auctionID) == "" {
		return fmt.Errorf("missing seller or auction id")
	}
	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(endAtRFC3339))
	if err != nil {
		return fmt.Errorf("invalid end_at")
	}
	if !endAt.After(time.Now()) {
		return fmt.Errorf("end_at must be in the future")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	a, err := s.repo.LockAuctionBySellerForUpdate(ctx, tx, auctionID, sellerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("auction not found")
		}
		return err
	}
	if a.Status != "closed" || a.TotalBids != 0 || strings.TrimSpace(a.WinnerID) != "" {
		return ErrAuctionReopenNotAllowed
	}
	nBids, err := s.repo.CountAuctionBidsTx(ctx, tx, auctionID)
	if err != nil {
		return err
	}
	if nBids > 0 {
		return ErrAuctionReopenNotAllowed
	}
	nHeld, err := s.repo.CountHeldBidHoldsTx(ctx, tx, auctionID)
	if err != nil {
		return err
	}
	if nHeld > 0 {
		return ErrAuctionReopenNotAllowed
	}

	ok, balBefore, balAfter, err := s.userRepo.DeductListingDepositTx(ctx, tx, sellerID, a.StartPrice)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("insufficient credit for start price (%d THB required)", a.StartPrice)
	}

	n, err := s.repo.ApplyAuctionReopenTx(ctx, tx, auctionID, sellerID, endAt)
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrAuctionReopenNotAllowed
	}
	if err := s.repo.InsertListingDepositHoldTx(ctx, tx, sellerID, auctionID, a.StartPrice, balBefore, balAfter, "หักมัดจำประกาศเมื่อเปิดประมูลรอบใหม่"); err != nil {
		return err
	}
	return tx.Commit()
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

func generateAuctionID() string {
	return fmt.Sprintf("AUC-%s", strings.ToUpper(strings.ReplaceAll(time.Now().Format("20060102-150405.000"), ".", "")))
}
