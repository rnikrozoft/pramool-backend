package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type dsarRequestRow struct {
	ID          int64     `bun:"id"`
	UserID      string    `bun:"user_id"`
	RequestType string    `bun:"request_type"`
	Status      string    `bun:"status"`
	CreatedAt   time.Time `bun:"created_at"`
}

func (r privacy) GetDSARRequestForUser(ctx context.Context, dsarID int64, userID string) (*dsarRequestRow, error) {
	row := new(dsarRequestRow)
	err := r.bun.NewSelect().
		TableExpr("dsar_requests").
		Column("id", "user_id", "request_type", "status", "created_at").
		Where("id = ? AND user_id = ?", dsarID, strings.TrimSpace(userID)).
		Scan(ctx, row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPrivacyNotFound
	}
	return row, err
}

func (r privacy) SaveDSARExport(ctx context.Context, dsarID int64, payload []byte) error {
	_, err := r.bun.NewRaw(`
INSERT INTO dsar_request_exports (dsar_id, export_json)
VALUES (?, ?::jsonb)
ON CONFLICT (dsar_id) DO UPDATE SET export_json = EXCLUDED.export_json, created_at = NOW()
`, dsarID, string(payload)).Exec(ctx)
	return err
}

func (r privacy) GetDSARExportJSON(ctx context.Context, dsarID int64, userID string) ([]byte, error) {
	var raw json.RawMessage
	err := r.bun.NewRaw(`
SELECT e.export_json
FROM dsar_request_exports e
JOIN dsar_requests d ON d.id = e.dsar_id
WHERE e.dsar_id = ? AND d.user_id = ?
`, dsarID, strings.TrimSpace(userID)).Scan(ctx, &raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPrivacyNotFound
	}
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r privacy) CompleteDSARRequest(ctx context.Context, dsarID int64) error {
	res, err := r.bun.NewRaw(`
UPDATE dsar_requests
SET status = 'completed', completed_at = NOW(), updated_at = NOW()
WHERE id = ?
`, dsarID).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPrivacyNotFound
	}
	return nil
}

type exportProfileRow struct {
	UserID         string    `json:"user_id"`
	Tel            string    `json:"tel"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Facebook       string    `json:"facebook"`
	AddressPrimary string    `json:"address_primary"`
	Address        string    `json:"address"`
	Soi            string    `json:"soi"`
	Road           string    `json:"road"`
	SubDistrict    string    `json:"sub_district"`
	District       string    `json:"district"`
	Province       string    `json:"province"`
	ZipCode        string    `json:"zip_code"`
	Credit         int64     `json:"credit"`
	MarketingOptIn bool      `json:"marketing_opt_in"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type exportConsentRow struct {
	ConsentType   string    `json:"consent_type"`
	PolicyVersion string    `json:"policy_version"`
	CreatedAt     time.Time `json:"created_at"`
}

type exportDSARRow struct {
	ID          int64      `json:"id"`
	RequestType string     `json:"request_type"`
	Status      string     `json:"status"`
	UserNote    *string    `json:"user_note,omitempty"`
	AdminNote   *string    `json:"admin_note,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"`
}

type exportNotificationRow struct {
	ID        int64      `json:"id"`
	Kind      string     `json:"kind"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type exportSellerAuctionRow struct {
	AuctionID  string    `json:"auction_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	StartPrice int64     `json:"start_price"`
	CurrentBid int64     `json:"current_bid"`
	EndAt      time.Time `json:"end_at"`
	WinnerID   *string   `json:"winner_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type exportBidAuctionRow struct {
	AuctionID  string    `json:"auction_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	CurrentBid int64     `json:"current_bid"`
	EndAt      time.Time `json:"end_at"`
	LastBidAt  time.Time `json:"last_bid_at"`
}

type exportWalletTxRow struct {
	TransactionID int64     `json:"transaction_id"`
	ChargeID      string    `json:"charge_id"`
	Amount        int64     `json:"amount"`
	FeeAmount     int64     `json:"fee_amount"`
	CreditAmount  int64     `json:"credit_amount"`
	Status        string    `json:"status"`
	Paid          bool      `json:"paid"`
	Credited      bool      `json:"credited"`
	CreatedAt     time.Time `json:"created_at"`
}

type exportWithdrawalRow struct {
	WithdrawalID   int64     `json:"withdrawal_id"`
	Amount         int64     `json:"amount"`
	FeeAmount      int64     `json:"fee_amount"`
	TransferAmount int64     `json:"transfer_amount"`
	Status         string    `json:"status"`
	FailureReason  *string   `json:"failure_reason,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type exportReviewRow struct {
	AuctionID string    `json:"auction_id"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

func (r privacy) BuildUserDataExport(ctx context.Context, userID string) (map[string]any, error) {
	userID = strings.TrimSpace(userID)
	out := map[string]any{
		"exported_at": time.Now().UTC().Format(time.RFC3339),
		"user_id":     userID,
	}

	var profile exportProfileRow
	err := r.bun.NewRaw(`
SELECT user_id::text, tel,
       COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(email, ''), COALESCE(facebook, ''),
       COALESCE(address_primary, ''), COALESCE(address, ''), COALESCE(soi, ''), COALESCE(road, ''),
       COALESCE(sub_district, ''), COALESCE(district, ''), COALESCE(province, ''), COALESCE(zip_code, ''),
       COALESCE(credit, 0), COALESCE(marketing_opt_in, FALSE), created_at, updated_at
FROM users WHERE user_id = ?
`, userID).Scan(ctx, &profile)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPrivacyNotFound
	}
	if err != nil {
		return nil, err
	}
	out["profile"] = profile

	var consents []exportConsentRow
	_ = r.bun.NewRaw(`
SELECT consent_type, policy_version, created_at
FROM consent_log WHERE user_id = ? ORDER BY created_at DESC LIMIT 200
`, userID).Scan(ctx, &consents)
	out["consents"] = consents

	var dsarRows []exportDSARRow
	_ = r.bun.NewRaw(`
SELECT id, request_type, status, user_note, admin_note, created_at, completed_at, due_at
FROM dsar_requests WHERE user_id = ? ORDER BY created_at DESC
`, userID).Scan(ctx, &dsarRows)
	out["dsar_requests"] = dsarRows

	var notifications []exportNotificationRow
	_ = r.bun.NewRaw(`
SELECT id, kind, title, body, read_at, expires_at, created_at
FROM user_notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT 500
`, userID).Scan(ctx, &notifications)
	out["notifications"] = notifications

	var sellerAuctions []exportSellerAuctionRow
	_ = r.bun.NewRaw(`
SELECT auction_id, title, status, start_price, current_bid, end_at, winner_id::text, created_at
FROM auctions WHERE seller_id = ? ORDER BY created_at DESC LIMIT 500
`, userID).Scan(ctx, &sellerAuctions)
	out["auctions_as_seller"] = sellerAuctions

	var bidAuctions []exportBidAuctionRow
	_ = r.bun.NewRaw(`
SELECT a.auction_id, a.title, a.status, a.current_bid, a.end_at, abp.last_bid_at
FROM auction_bid_participants abp
JOIN auctions a ON a.auction_id = abp.auction_id
WHERE abp.bidder_user_id = ?
ORDER BY abp.last_bid_at DESC LIMIT 500
`, userID).Scan(ctx, &bidAuctions)
	out["auctions_as_bidder"] = bidAuctions

	var walletTx []exportWalletTxRow
	_ = r.bun.NewRaw(`
SELECT transaction_id, charge_id, amount, fee_amount, credit_amount, status, paid, credited, created_at
FROM transactions WHERE user_id = ? ORDER BY created_at DESC LIMIT 500
`, userID).Scan(ctx, &walletTx)
	out["wallet_transactions"] = walletTx

	var withdrawals []exportWithdrawalRow
	_ = r.bun.NewRaw(`
SELECT withdrawal_id, amount, fee_amount, transfer_amount, status, failure_reason, created_at
FROM withdrawals WHERE user_id = ? ORDER BY created_at DESC LIMIT 200
`, userID).Scan(ctx, &withdrawals)
	out["withdrawals"] = withdrawals

	var reviews []exportReviewRow
	_ = r.bun.NewRaw(`
SELECT auction_id, rating::float8, created_at FROM auction_seller_reviews
WHERE buyer_user_id = ? OR seller_id = ?
ORDER BY created_at DESC LIMIT 200
`, userID, userID).Scan(ctx, &reviews)
	out["reviews"] = reviews

	return out, nil
}

func (r privacy) PurgeUserPersonalDataOnAnonymize(ctx context.Context, userID string) error {
	userID = strings.TrimSpace(userID)
	if _, err := r.bun.NewRaw(`DELETE FROM user_notifications WHERE user_id = ?`, userID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.bun.NewRaw(`UPDATE consent_log SET user_id = NULL, tel = NULL WHERE user_id = ?`, userID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.bun.NewRaw(`UPDATE dsar_requests SET user_note = NULL WHERE user_id = ?`, userID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.bun.NewRaw(`UPDATE data_disclosure_log SET ip_address = NULL WHERE actor_user_id = ? OR subject_user_id = ?`, userID, userID).Exec(ctx); err != nil {
		return err
	}
	if _, err := r.bun.NewRaw(`UPDATE user_restriction_appeals SET reason = '[ลบแล้ว]' WHERE user_id = ?`, userID).Exec(ctx); err != nil {
		return err
	}
	_, err := r.bun.NewRaw(`
DELETE FROM dsar_request_exports e
USING dsar_requests d
WHERE d.id = e.dsar_id AND d.user_id = ?
`, userID).Exec(ctx)
	return err
}
