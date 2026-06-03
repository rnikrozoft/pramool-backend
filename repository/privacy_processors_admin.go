package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

func (r privacy) GetNationalIDEnc(ctx context.Context, userID string) (string, error) {
	var enc sql.NullString
	err := r.bun.NewRaw(`SELECT national_id_enc FROM users WHERE user_id = ?`, strings.TrimSpace(userID)).Scan(ctx, &enc)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrPrivacyNotFound
	}
	if err != nil {
		return "", err
	}
	if !enc.Valid {
		return "", nil
	}
	return enc.String, nil
}

func (r privacy) ListAllDataProcessors(ctx context.Context) ([]DataProcessorRow, error) {
	var rows []DataProcessorRow
	err := r.bun.NewSelect().
		TableExpr("data_processors").
		Column("processor_id", "name", "purpose", "data_categories", "location", "privacy_url", "dpa_status", "is_active", "sort_order").
		OrderExpr("sort_order ASC, processor_id ASC").
		Scan(ctx, &rows)
	return rows, err
}

func (r privacy) CreateDataProcessor(ctx context.Context, row DataProcessorRow) (*DataProcessorRow, error) {
	_, err := r.bun.NewInsert().Model(&row).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r privacy) UpdateDataProcessor(ctx context.Context, id int, row DataProcessorRow) error {
	res, err := r.bun.NewUpdate().TableExpr("data_processors").
		Set("name = ?", row.Name).
		Set("purpose = ?", row.Purpose).
		Set("data_categories = ?", row.DataCategories).
		Set("location = ?", row.Location).
		Set("privacy_url = ?", row.PrivacyURL).
		Set("dpa_status = ?", row.DPAStatus).
		Set("is_active = ?", row.IsActive).
		Set("sort_order = ?", row.SortOrder).
		Set("updated_at = NOW()").
		Where("processor_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPrivacyNotFound
	}
	return nil
}

func (r privacy) DeleteDataProcessor(ctx context.Context, id int) error {
	res, err := r.bun.NewDelete().TableExpr("data_processors").Where("processor_id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPrivacyNotFound
	}
	return nil
}

func (r privacy) DSARExportExists(ctx context.Context, dsarID int64) (bool, error) {
	var exists bool
	err := r.bun.NewRaw(`SELECT EXISTS(SELECT 1 FROM dsar_request_exports WHERE dsar_id = ?)`, dsarID).Scan(ctx, &exists)
	return exists, err
}
