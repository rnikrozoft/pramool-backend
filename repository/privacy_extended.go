package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type AccountDeletionBlockers struct {
	CreditBalance        int64
	ActiveSellerAuctions int64
	ActiveBidAuctions    int64
	PendingSellerShip    int64
	PendingBuyerConfirm  int64
	PendingWithdrawals   int64
	AlreadyDeleted       bool
}

type DataProcessorRow struct {
	ProcessorID    int     `bun:"processor_id,pk,autoincrement"`
	Name           string  `bun:"name,notnull"`
	Purpose        string  `bun:"purpose,notnull"`
	DataCategories string  `bun:"data_categories,notnull"`
	Location       string  `bun:"location,notnull"`
	PrivacyURL     *string `bun:"privacy_url"`
	DPAStatus      string  `bun:"dpa_status,notnull"`
	IsActive       bool    `bun:"is_active,notnull"`
	SortOrder      int     `bun:"sort_order,notnull"`
}

type RetentionJobDefinitionRow struct {
	JobID            string     `bun:"job_id,pk"`
	NameTH           string     `bun:"name_th,notnull"`
	Description      string     `bun:"description,notnull"`
	RetentionDays    int        `bun:"retention_days,notnull"`
	BatchRunner      string     `bun:"batch_runner,notnull"`
	IsEnabled        bool       `bun:"is_enabled,notnull"`
	LastRunAt        *time.Time `bun:"last_run_at"`
	LastDeletedCount *int64     `bun:"last_deleted_count"`
}

func (r privacy) GetAccountDeletionBlockers(ctx context.Context, userID string) (*AccountDeletionBlockers, error) {
	userID = strings.TrimSpace(userID)
	out := &AccountDeletionBlockers{}
	err := r.bun.NewRaw(`
SELECT
  COALESCE(u.credit, 0),
  (SELECT COUNT(*)::bigint FROM auctions a WHERE a.seller_id = u.user_id AND a.status = 'active' AND a.end_at > NOW()),
  (SELECT COUNT(DISTINCT abp.auction_id)::bigint FROM auction_bid_participants abp
     JOIN auctions a ON a.auction_id = abp.auction_id
     WHERE abp.bidder_user_id = u.user_id AND a.status = 'active' AND a.end_at > NOW()),
  (SELECT COUNT(*)::bigint FROM auctions a
     WHERE a.seller_id = u.user_id AND a.status = 'closed' AND a.seller_payout_at IS NULL
       AND a.winner_id IS NOT NULL AND a.seller_shipped_at IS NULL
       AND a.buyer_escrow_refunded_at IS NULL),
  (SELECT COUNT(*)::bigint FROM auctions a
     WHERE a.winner_id = u.user_id AND a.status = 'closed' AND a.seller_payout_at IS NULL
       AND a.winner_id IS NOT NULL AND a.buyer_received_at IS NULL),
  (SELECT COUNT(*)::bigint FROM withdrawals w
     WHERE w.user_id = u.user_id AND w.status IN ('pending', 'processing')),
  (u.deleted_at IS NOT NULL)
FROM users u
WHERE u.user_id = ?
`, userID).Scan(ctx,
		&out.CreditBalance,
		&out.ActiveSellerAuctions,
		&out.ActiveBidAuctions,
		&out.PendingSellerShip,
		&out.PendingBuyerConfirm,
		&out.PendingWithdrawals,
		&out.AlreadyDeleted,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPrivacyNotFound
		}
		return nil, err
	}
	return out, nil
}

func (r privacy) AnonymizeUserAccount(ctx context.Context, userID string) error {
	userID = strings.TrimSpace(userID)
	newTel, err := randomDeletedTel()
	if err != nil {
		return err
	}
	res, err := r.bun.NewRaw(`
UPDATE users SET
  first_name = '[ลบแล้ว]',
  last_name = '-',
  email = '',
  facebook = '',
  address_primary = '',
  address = '',
  soi = '',
  road = '',
  sub_district = '',
  district = '',
  province = '',
  zip_code = '',
  bank_id = NULL,
  bank_account_name = NULL,
  bank_account_number = NULL,
  omise_recipient_id = NULL,
  password_hash = NULL,
  marketing_opt_in = FALSE,
  national_id_hash = NULL,
  national_id_enc = NULL,
  tel = ?,
  deleted_at = NOW(),
  anonymized_at = NOW(),
  suspended_at = COALESCE(suspended_at, NOW()),
  suspended_reason = COALESCE(NULLIF(TRIM(suspended_reason), ''), 'บัญชีถูกลบตามคำขอ PDPA'),
  updated_at = NOW()
WHERE user_id = ? AND deleted_at IS NULL
`, newTel, userID).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		var exists bool
		_ = r.bun.NewRaw(`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)`, userID).Scan(ctx, &exists)
		if !exists {
			return ErrPrivacyNotFound
		}
		return errors.New("account already deleted")
	}
	return r.PurgeUserPersonalDataOnAnonymize(ctx, userID)
}

func randomDeletedTel() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("9%09d", n.Int64()), nil
}

func (r privacy) IsAccountDeleted(ctx context.Context, userID string) (bool, error) {
	var deleted bool
	err := r.bun.NewRaw(`SELECT deleted_at IS NOT NULL FROM users WHERE user_id = ?`, strings.TrimSpace(userID)).Scan(ctx, &deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return deleted, err
}

func (r privacy) IsAccountDeletedByTel(ctx context.Context, tel string) (bool, error) {
	var deleted bool
	err := r.bun.NewRaw(`SELECT deleted_at IS NOT NULL FROM users WHERE tel = ?`, strings.TrimSpace(tel)).Scan(ctx, &deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return deleted, err
}

func (r privacy) LogDataDisclosure(ctx context.Context, disclosureType, actorUserID, subjectUserID, referenceID, ip string) error {
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	var refPtr *string
	if referenceID != "" {
		refPtr = &referenceID
	}
	_, err := r.bun.NewRaw(`
INSERT INTO data_disclosure_log (disclosure_type, actor_user_id, subject_user_id, reference_id, ip_address)
VALUES (?, ?, ?, ?, ?)
`, disclosureType, actorUserID, subjectUserID, refPtr, ipPtr).Exec(ctx)
	return err
}

func (r privacy) ListActiveDataProcessors(ctx context.Context) ([]DataProcessorRow, error) {
	var rows []DataProcessorRow
	err := r.bun.NewSelect().
		TableExpr("data_processors").
		Column("processor_id", "name", "purpose", "data_categories", "location", "privacy_url", "dpa_status", "is_active", "sort_order").
		Where("is_active = TRUE").
		OrderExpr("sort_order ASC, processor_id ASC").
		Scan(ctx, &rows)
	return rows, err
}

func (r privacy) ListRetentionJobDefinitions(ctx context.Context) ([]RetentionJobDefinitionRow, error) {
	var rows []RetentionJobDefinitionRow
	err := r.bun.NewSelect().
		TableExpr("retention_job_definitions").
		Column("job_id", "name_th", "description", "retention_days", "batch_runner", "is_enabled", "last_run_at", "last_deleted_count").
		OrderExpr("job_id ASC").
		Scan(ctx, &rows)
	return rows, err
}

func (r privacy) TouchRetentionJobRun(ctx context.Context, jobID string, deleted int64) error {
	_, err := r.bun.NewRaw(`
UPDATE retention_job_definitions
SET last_run_at = NOW(), last_deleted_count = ?, updated_at = NOW()
WHERE job_id = ?
`, deleted, jobID).Exec(ctx)
	return err
}

func (r privacy) SetMarketingOptIn(ctx context.Context, userID string, optIn bool) error {
	res, err := r.bun.NewRaw(`
UPDATE users SET marketing_opt_in = ?, updated_at = NOW()
WHERE user_id = ? AND deleted_at IS NULL
`, optIn, strings.TrimSpace(userID)).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPrivacyNotFound
	}
	return nil
}

func (r privacy) GetMarketingOptIn(ctx context.Context, userID string) (bool, error) {
	var optIn bool
	err := r.bun.NewRaw(`SELECT marketing_opt_in FROM users WHERE user_id = ?`, strings.TrimSpace(userID)).Scan(ctx, &optIn)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrPrivacyNotFound
	}
	return optIn, err
}

func (r privacy) MarkDSARDeletionExecuted(ctx context.Context, dsarID int64) error {
	_, err := r.bun.NewUpdate().TableExpr("dsar_requests").
		Set("deletion_executed_at = NOW()").
		Set("updated_at = NOW()").
		Where("id = ?", dsarID).
		Exec(ctx)
	return err
}
