package repository

import (
	"context"
	"errors"
	"time"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

var ErrPrivacyNotFound = errors.New("privacy record not found")

type Privacy interface {
	InsertConsentLog(ctx context.Context, row entity.ConsentLog) error
	CreateDSARRequest(ctx context.Context, row entity.DSARRequest) (*entity.DSARRequest, error)
	ListDSARRequestsByUser(ctx context.Context, userID string, limit int) ([]entity.DSARRequest, error)
	PurgeStaleTelVerify(ctx context.Context, retentionDays int) (int64, error)
	GetAccountDeletionBlockers(ctx context.Context, userID string) (*AccountDeletionBlockers, error)
	AnonymizeUserAccount(ctx context.Context, userID string) error
	IsAccountDeleted(ctx context.Context, userID string) (bool, error)
	IsAccountDeletedByTel(ctx context.Context, tel string) (bool, error)
	LogDataDisclosure(ctx context.Context, disclosureType, actorUserID, subjectUserID, referenceID, ip string) error
	ListActiveDataProcessors(ctx context.Context) ([]DataProcessorRow, error)
	ListRetentionJobDefinitions(ctx context.Context) ([]RetentionJobDefinitionRow, error)
	TouchRetentionJobRun(ctx context.Context, jobID string, deleted int64) error
	SetMarketingOptIn(ctx context.Context, userID string, optIn bool) error
	GetMarketingOptIn(ctx context.Context, userID string) (bool, error)
	MarkDSARDeletionExecuted(ctx context.Context, dsarID int64) error
	BuildUserDataExport(ctx context.Context, userID string) (map[string]any, error)
	SaveDSARExport(ctx context.Context, dsarID int64, payload []byte) error
	GetDSARExportJSON(ctx context.Context, dsarID int64, userID string) ([]byte, error)
	CompleteDSARRequest(ctx context.Context, dsarID int64) error
	GetDSARRequestForUser(ctx context.Context, dsarID int64, userID string) (*dsarRequestRow, error)
	PurgeUserPersonalDataOnAnonymize(ctx context.Context, userID string) error
	ListAllDataProcessors(ctx context.Context) ([]DataProcessorRow, error)
	CreateDataProcessor(ctx context.Context, row DataProcessorRow) (*DataProcessorRow, error)
	UpdateDataProcessor(ctx context.Context, id int, row DataProcessorRow) error
	DeleteDataProcessor(ctx context.Context, id int) error
	GetNationalIDEnc(ctx context.Context, userID string) (string, error)
	DSARExportExists(ctx context.Context, dsarID int64) (bool, error)
}

type privacy struct {
	bun *bun.DB
}

func NewPrivacyRepository(bun *bun.DB) Privacy {
	return privacy{bun: bun}
}

func (r privacy) InsertConsentLog(ctx context.Context, row entity.ConsentLog) error {
	_, err := r.bun.NewInsert().Model(&row).Exec(ctx)
	return err
}

func (r privacy) CreateDSARRequest(ctx context.Context, row entity.DSARRequest) (*entity.DSARRequest, error) {
	due := time.Now().Add(30 * 24 * time.Hour)
	row.DueAt = &due
	_, err := r.bun.NewInsert().Model(&row).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r privacy) ListDSARRequestsByUser(ctx context.Context, userID string, limit int) ([]entity.DSARRequest, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows []entity.DSARRequest
	err := r.bun.NewSelect().
		Model(&rows).
		Where("user_id = ?", userID).
		OrderExpr("created_at DESC").
		Limit(limit).
		Scan(ctx)
	return rows, err
}

func (r privacy) PurgeStaleTelVerify(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	res, err := r.bun.NewRaw(`
DELETE FROM tel_verify tv
WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.tel = tv.tel)
  AND tv.updated_at < ?
`, cutoff).Exec(ctx)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
