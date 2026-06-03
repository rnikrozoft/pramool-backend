package repository

import (
	"context"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type ShipmentEventRow struct {
	Status     string
	Location   string
	Note       string
	OccurredAt time.Time
}

type WinnerShippingAddressRow struct {
	WinnerUserID   string
	FirstName      string
	LastName       string
	Phone          string
	AddressPrimary string
	Address        string
	Soi            string
	Road           string
	SubDistrict    string
	District       string
	Province       string
	ZipCode        string
}

type ShipmentRepository interface {
	MarkSellerShippedWithTracking(ctx context.Context, auctionID, sellerID, carrierCode, carrierName, trackingNumber, shipmentStatus string) (int64, error)
	UpdateAuctionShipmentStatus(ctx context.Context, auctionID, status string) error
	UpdateAuctionShipmentCarrier(ctx context.Context, auctionID, carrierCode, carrierName string) error
	ReplaceAuctionShipmentEvents(ctx context.Context, auctionID string, events []ShipmentEventRow) error
	ListAuctionShipmentEvents(ctx context.Context, auctionID string) ([]ShipmentEventRow, error)
	GetAuctionShipmentAccess(ctx context.Context, auctionID string) (sellerID, winnerID, carrierCode, carrierName, trackingNumber, shipmentStatus string, sellerShipped bool, err error)
	GetWinnerShippingAddressForSeller(ctx context.Context, auctionID, sellerUserID string) (WinnerShippingAddressRow, error)
}

type shipmentRepo struct {
	bun *bun.DB
}

func NewShipmentRepository(db *bun.DB) ShipmentRepository {
	return shipmentRepo{bun: db}
}

func (r shipmentRepo) MarkSellerShippedWithTracking(
	ctx context.Context,
	auctionID, sellerID, carrierCode, carrierName, trackingNumber, shipmentStatus string,
) (int64, error) {
	res, err := r.bun.ExecContext(ctx, `
		UPDATE auctions
		SET seller_shipped_at = NOW(),
		    carrier_code = ?,
		    carrier_name = ?,
		    tracking_number = ?,
		    shipment_status = ?,
		    updated_at = NOW()
		WHERE auction_id = ?
		  AND seller_id = ?
		  AND status = 'closed'
		  AND winner_id IS NOT NULL
		  AND seller_shipped_at IS NULL
		  AND seller_payout_at IS NULL
	`, carrierCode, carrierName, trackingNumber, shipmentStatus, auctionID, sellerID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r shipmentRepo) UpdateAuctionShipmentStatus(ctx context.Context, auctionID, status string) error {
	_, err := r.bun.ExecContext(ctx, `
		UPDATE auctions
		SET shipment_status = ?, updated_at = NOW()
		WHERE auction_id = ?
	`, strings.ToLower(strings.TrimSpace(status)), auctionID)
	return err
}

func (r shipmentRepo) UpdateAuctionShipmentCarrier(ctx context.Context, auctionID, carrierCode, carrierName string) error {
	_, err := r.bun.ExecContext(ctx, `
		UPDATE auctions
		SET carrier_code = ?,
		    carrier_name = ?,
		    updated_at = NOW()
		WHERE auction_id = ?
	`, strings.TrimSpace(carrierCode), strings.TrimSpace(carrierName), auctionID)
	return err
}

func (r shipmentRepo) ReplaceAuctionShipmentEvents(ctx context.Context, auctionID string, events []ShipmentEventRow) error {
	tx, err := r.bun.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM auction_shipment_events WHERE auction_id = ?`, auctionID); err != nil {
		return err
	}
	for _, ev := range events {
		at := ev.OccurredAt
		if at.IsZero() {
			at = time.Now()
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO auction_shipment_events (auction_id, status, location, note, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, auctionID, ev.Status, shipmentNullIfEmpty(ev.Location), shipmentNullIfEmpty(ev.Note), at); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r shipmentRepo) ListAuctionShipmentEvents(ctx context.Context, auctionID string) ([]ShipmentEventRow, error) {
	rows, err := r.bun.QueryContext(ctx, `
		SELECT status, COALESCE(location, ''), COALESCE(note, ''), created_at
		FROM auction_shipment_events
		WHERE auction_id = ?
		ORDER BY created_at DESC, event_id DESC
	`, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ShipmentEventRow, 0)
	for rows.Next() {
		var row ShipmentEventRow
		if err := rows.Scan(&row.Status, &row.Location, &row.Note, &row.OccurredAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r shipmentRepo) GetAuctionShipmentAccess(ctx context.Context, auctionID string) (sellerID, winnerID, carrierCode, carrierName, trackingNumber, shipmentStatus string, sellerShipped bool, err error) {
	err = r.bun.NewRaw(`
		SELECT seller_id,
		       COALESCE(a.winner_id::text, ''),
		       COALESCE(carrier_code, ''),
		       COALESCE(carrier_name, ''),
		       COALESCE(tracking_number, ''),
		       COALESCE(shipment_status, 'pending'),
		       (seller_shipped_at IS NOT NULL)
		FROM auctions
		WHERE auction_id = ?
		LIMIT 1
	`, auctionID).Scan(ctx, &sellerID, &winnerID, &carrierCode, &carrierName, &trackingNumber, &shipmentStatus, &sellerShipped)
	return
}

func (r shipmentRepo) GetWinnerShippingAddressForSeller(ctx context.Context, auctionID, sellerUserID string) (WinnerShippingAddressRow, error) {
	var row WinnerShippingAddressRow
	err := r.bun.NewRaw(`
		SELECT COALESCE(a.winner_id::text, ''),
		       COALESCE(u.first_name, ''),
		       COALESCE(u.last_name, ''),
		       COALESCE(u.tel, ''),
		       COALESCE(u.address_primary, ''),
		       COALESCE(u.address, ''),
		       COALESCE(u.soi, ''),
		       COALESCE(u.road, ''),
		       COALESCE(u.sub_district, ''),
		       COALESCE(u.district, ''),
		       COALESCE(u.province, ''),
		       COALESCE(u.zip_code, '')
		FROM auctions a
		INNER JOIN users u ON u.user_id = a.winner_id
		WHERE a.auction_id = ?
		  AND a.seller_id = ?
		  AND a.status = 'closed'
		  AND a.winner_id IS NOT NULL
		  AND a.seller_shipped_at IS NULL
		  AND a.seller_payout_at IS NULL
		LIMIT 1
	`, auctionID, sellerUserID).Scan(
		ctx,
		&row.WinnerUserID,
		&row.FirstName,
		&row.LastName,
		&row.Phone,
		&row.AddressPrimary,
		&row.Address,
		&row.Soi,
		&row.Road,
		&row.SubDistrict,
		&row.District,
		&row.Province,
		&row.ZipCode,
	)
	return row, err
}

func shipmentNullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
