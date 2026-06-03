package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/internal/tracking"
	"github.com/rnikrozoft/pramool-core/internal/tracking/trackingmore"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/repository"
)

type ShipmentService interface {
	ListCarriers() []dto.ShipmentCarrierItem
	MarkSellerShipped(ctx context.Context, auctionID, sellerUserID, carrierCode, trackingNumber string) error
	GetWinnerShippingAddress(ctx context.Context, auctionID, sellerUserID, clientIP string) (*dto.WinnerShippingAddressResponse, error)
	RefreshTracking(ctx context.Context, auctionID, userID string) (*dto.ShipmentTrackingResponse, error)
}

type shipmentSvc struct {
	repo     repository.ShipmentRepository
	tracking *trackingmore.Client
	privacy  PrivacyService
}

func NewShipmentService(repo repository.ShipmentRepository, trackingAPIKey string, privacy PrivacyService) ShipmentService {
	return shipmentSvc{
		repo:     repo,
		tracking: trackingmore.NewClient(trackingAPIKey),
		privacy:  privacy,
	}
}

func (s shipmentSvc) ListCarriers() []dto.ShipmentCarrierItem {
	carriers := tracking.ListCarriers()
	out := make([]dto.ShipmentCarrierItem, 0, len(carriers))
	for _, c := range carriers {
		out = append(out, dto.ShipmentCarrierItem{Code: c.Code, Label: c.Label})
	}
	return out
}

func (s shipmentSvc) GetWinnerShippingAddress(ctx context.Context, auctionID, sellerUserID, clientIP string) (*dto.WinnerShippingAddressResponse, error) {
	aid := strings.TrimSpace(auctionID)
	sid := strings.TrimSpace(sellerUserID)
	if aid == "" || sid == "" {
		return nil, ErrMarkShippedNotAllowed
	}
	row, err := s.repo.GetWinnerShippingAddressForSeller(ctx, aid, sid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMarkShippedNotAllowed
		}
		return nil, err
	}
	phone := strings.TrimSpace(row.Phone)
	if s.privacy != nil && strings.TrimSpace(row.WinnerUserID) != "" {
		_ = s.privacy.LogWinnerAddressDisclosure(ctx, sid, strings.TrimSpace(row.WinnerUserID), aid, clientIP)
	}
	return &dto.WinnerShippingAddressResponse{
		FirstName:         strings.TrimSpace(row.FirstName),
		LastName:          strings.TrimSpace(row.LastName),
		Phone:             phone,
		Address:           formatWinnerShippingAddress(row),
		DisclosureNotice:  "ข้อมูลนี้ใช้เพื่อจัดส่งสินค้าตามสัญญาประมูลเท่านั้น ห้ามใช้นอกวัตถุประสงค์ การเข้าถึงถูกบันทึกในระบบ",
		PhoneFullyVisible: true,
	}, nil
}

func formatWinnerShippingAddress(row repository.WinnerShippingAddressRow) string {
	parts := make([]string, 0, 8)
	for _, p := range []string{row.AddressPrimary, row.Address, row.Soi, row.Road, row.SubDistrict, row.District, row.Province} {
		if s := strings.TrimSpace(p); s != "" {
			parts = append(parts, s)
		}
	}
	if zip := strings.TrimSpace(row.ZipCode); zip != "" {
		parts = append(parts, zip)
	}
	return strings.Join(parts, " ")
}

func (s shipmentSvc) MarkSellerShipped(ctx context.Context, auctionID, sellerUserID, carrierCode, trackingNumber string) error {
	if s.tracking == nil || !s.tracking.Enabled() {
		return fmt.Errorf("ระบบติดตามพัสดุยังไม่พร้อม กรุณาติดต่อผู้ดูแล")
	}
	aid := strings.TrimSpace(auctionID)
	sid := strings.TrimSpace(sellerUserID)
	carrierCode = strings.TrimSpace(carrierCode)
	trackingNumber = strings.TrimSpace(trackingNumber)
	if aid == "" || sid == "" {
		return ErrMarkShippedNotAllowed
	}
	if trackingNumber == "" || len(trackingNumber) > 64 {
		return ErrInvalidTrackingNumber
	}

	resolvedCode, resolvedLabel, err := s.resolveCarrier(ctx, carrierCode, trackingNumber)
	if err != nil {
		return err
	}

	result, err := s.tracking.CreateTracking(ctx, resolvedCode, trackingNumber)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(strings.ToLower(err.Error()), "timeout") {
			return fmt.Errorf("ติดต่อ TrackingMore ไม่สำเร็จ กรุณาลองใหม่")
		}
		return err
	}
	status := strings.ToLower(strings.TrimSpace(result.Status))
	if status == "" {
		status = "pending"
	}

	n, err := s.repo.MarkSellerShippedWithTracking(ctx, aid, sid, resolvedCode, resolvedLabel, trackingNumber, status)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrMarkShippedNotAllowed
	}
	if len(result.Events) > 0 {
		_ = s.repo.ReplaceAuctionShipmentEvents(ctx, aid, mapTrackingEvents(result.Events))
	}
	return nil
}

func (s shipmentSvc) RefreshTracking(ctx context.Context, auctionID, userID string) (*dto.ShipmentTrackingResponse, error) {
	if s.tracking == nil || !s.tracking.Enabled() {
		return nil, fmt.Errorf("ระบบติดตามพัสดุยังไม่พร้อม กรุณาติดต่อผู้ดูแล")
	}
	aid := strings.TrimSpace(auctionID)
	uid := strings.TrimSpace(userID)
	if aid == "" || uid == "" {
		return nil, ErrShipmentAccessDenied
	}

	sellerID, winnerID, carrierCode, carrierName, trackingNumber, shipmentStatus, sellerShipped, err := s.repo.GetAuctionShipmentAccess(ctx, aid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("auction not found")
		}
		return nil, err
	}
	if uid != strings.TrimSpace(sellerID) && uid != strings.TrimSpace(winnerID) {
		return nil, ErrShipmentAccessDenied
	}
	if !sellerShipped || carrierCode == "" || trackingNumber == "" {
		return nil, ErrShipmentNotRegistered
	}

	apiCarrierCode, apiCarrierName := s.trackingCarrierForFetch(carrierCode, carrierName, trackingNumber)
	if apiCarrierCode != carrierCode {
		carrierCode = apiCarrierCode
		carrierName = apiCarrierName
		_ = s.repo.UpdateAuctionShipmentCarrier(ctx, aid, apiCarrierCode, apiCarrierName)
	}

	var result *trackingmore.Result
	if tracking.IsTrackingMoreTestNumber(trackingNumber) {
		result, err = s.tracking.CreateTracking(ctx, apiCarrierCode, trackingNumber)
	} else {
		result, err = s.tracking.GetTracking(ctx, apiCarrierCode, trackingNumber)
	}
	if err != nil {
		return nil, err
	}
	status := strings.ToLower(strings.TrimSpace(result.Status))
	if status == "" {
		status = shipmentStatus
	}
	if status != shipmentStatus {
		if err := s.repo.UpdateAuctionShipmentStatus(ctx, aid, status); err != nil {
			return nil, err
		}
		shipmentStatus = status
	}
	if len(result.Events) > 0 {
		if err := s.repo.ReplaceAuctionShipmentEvents(ctx, aid, mapTrackingEvents(result.Events)); err != nil {
			return nil, err
		}
	}

	events, err := s.repo.ListAuctionShipmentEvents(ctx, aid)
	if err != nil {
		return nil, err
	}
	outEvents := make([]dto.ShipmentTrackingEvent, 0, len(events))
	for _, ev := range events {
		item := dto.ShipmentTrackingEvent{
			Status:   ev.Status,
			Location: ev.Location,
			Note:     ev.Note,
		}
		if !ev.OccurredAt.IsZero() {
			item.OccurredAt = ev.OccurredAt.Format(time.RFC3339)
		}
		outEvents = append(outEvents, item)
	}

	trackURL := ""
	if carrier, ok := tracking.CarrierByCode(carrierCode); ok && carrier.TrackURLFmt != "" {
		trackURL = fmt.Sprintf(carrier.TrackURLFmt, trackingNumber)
	}

	return &dto.ShipmentTrackingResponse{
		AuctionID:      aid,
		CarrierCode:    carrierCode,
		CarrierName:    carrierName,
		TrackingNumber: trackingNumber,
		ShipmentStatus: shipmentStatus,
		TrackURL:       trackURL,
		Events:         outEvents,
		CanConfirm:     shipmentStatus == "delivered",
	}, nil
}

func (s shipmentSvc) resolveCarrier(ctx context.Context, carrierCode, trackingNumber string) (code, label string, err error) {
	if tracking.IsTrackingMoreTestNumber(trackingNumber) {
		return tracking.TestCarrierCode, tracking.TestCarrierLabel, nil
	}

	carrierCode = strings.TrimSpace(carrierCode)
	if carrierCode != "" {
		if c, ok := tracking.CarrierByCode(carrierCode); ok {
			return c.Code, c.Label, nil
		}
		return carrierCode, tracking.CarrierLabel(carrierCode), nil
	}

	detected, err := s.tracking.DetectCarriers(ctx, trackingNumber)
	if err != nil {
		return "", "", err
	}
	known := tracking.KnownCarrierCodes()
	for _, row := range detected {
		if _, ok := known[strings.ToLower(row.Code)]; ok {
			return row.Code, tracking.CarrierLabel(row.Code), nil
		}
	}
	first := detected[0]
	label = first.Name
	if label == "" {
		label = tracking.CarrierLabel(first.Code)
	}
	return first.Code, label, nil
}

func (s shipmentSvc) trackingCarrierForFetch(carrierCode, carrierName, trackingNumber string) (code, label string) {
	if tracking.IsTrackingMoreTestNumber(trackingNumber) {
		return tracking.TestCarrierCode, tracking.TestCarrierLabel
	}
	return strings.TrimSpace(carrierCode), strings.TrimSpace(carrierName)
}

func mapTrackingEvents(in []trackingmore.Event) []repository.ShipmentEventRow {
	out := make([]repository.ShipmentEventRow, 0, len(in))
	for _, ev := range in {
		out = append(out, repository.ShipmentEventRow{
			Status:     ev.Status,
			Location:   ev.Location,
			Note:       ev.Note,
			OccurredAt: ev.OccurredAt,
		})
	}
	return out
}
