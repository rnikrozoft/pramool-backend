package repository

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

type BankRepository interface {
	ListActive(ctx context.Context) ([]entity.Bank, error)
	FindByID(ctx context.Context, bankID int64) (*entity.Bank, error)
}

type bankRepository struct {
	bun *bun.DB
}

func NewBankRepository(bun *bun.DB) BankRepository {
	return bankRepository{bun: bun}
}

func (r bankRepository) ListActive(ctx context.Context) ([]entity.Bank, error) {
	rows := make([]entity.Bank, 0)
	err := r.bun.NewRaw(`
		SELECT bank_id, bank_code, bank_name_th, bank_name_en, display_order
		FROM banks
		WHERE is_active = TRUE
		ORDER BY display_order ASC, bank_name_th ASC
	`).Scan(ctx, &rows)
	return rows, err
}

func (r bankRepository) FindByID(ctx context.Context, bankID int64) (*entity.Bank, error) {
	item := new(entity.Bank)
	err := r.bun.NewRaw(`
		SELECT bank_id, bank_code, bank_name_th, bank_name_en, display_order
		FROM banks
		WHERE bank_id = ? AND is_active = TRUE
	`, bankID).Scan(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
