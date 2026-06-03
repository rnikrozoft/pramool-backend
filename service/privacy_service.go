package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/internal/privacy"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

var (
	ErrInvalidPolicyVersion = errors.New("policy version is outdated or invalid")
	ErrInvalidDSARType      = errors.New("invalid dsar request type")
)

type PrivacyService interface {
	RecordConsents(ctx context.Context, userID, tel, ip, userAgent string, payload dto.ConsentPayload) error
	RecordCookieConsent(ctx context.Context, userID, tel, ip, userAgent string, req dto.CookieConsentRequest) error
	GetPolicyInfo() dto.PrivacyPolicyResponse
	CreateDSARRequest(ctx context.Context, userID string, req dto.CreateDSARRequest) (*dto.DSARRequestItem, error)
	ListMyDSARRequests(ctx context.Context, userID string) ([]dto.DSARRequestItem, error)
	PurgeStaleTelVerify(ctx context.Context) (int64, error)
	GetMarketingOptIn(ctx context.Context, userID string) (bool, error)
	UpdateMarketingConsent(ctx context.Context, userID, tel, ip, userAgent string, req dto.MarketingConsentRequest) error
	GetAccountDeletionReadiness(ctx context.Context, userID string) (*dto.AccountDeletionBlockersResponse, error)
	ExecuteAccountDeletion(ctx context.Context, userID string, dsarID int64) error
	GetDSARExport(ctx context.Context, userID string, dsarID int64) ([]byte, error)
	ListDataProcessors(ctx context.Context) ([]dto.DataProcessorItem, error)
	ListRetentionJobs(ctx context.Context) (*dto.RetentionJobListResponse, error)
	RunInlineRetentionJob(ctx context.Context, jobID string) (int64, error)
	LogWinnerAddressDisclosure(ctx context.Context, sellerUserID, buyerUserID, auctionID, ip string) error
	IsAccountDeleted(ctx context.Context, userID string) (bool, error)
	IsAccountDeletedByTel(ctx context.Context, tel string) (bool, error)
}

type privacyService struct {
	repo repository.Privacy
}

func NewPrivacyService(repo repository.Privacy) PrivacyService {
	return privacyService{repo: repo}
}

func (s privacyService) validateConsentPayload(payload dto.ConsentPayload) error {
	if payload.PrivacyPolicyVersion != privacy.PrivacyPolicyVersion {
		return ErrInvalidPolicyVersion
	}
	if payload.TermsVersion != privacy.TermsVersion {
		return ErrInvalidPolicyVersion
	}
	if !payload.AcceptPrivacy || !payload.AcceptTerms {
		return errors.New("consent required")
	}
	return nil
}

func (s privacyService) RecordConsents(ctx context.Context, userID, tel, ip, userAgent string, payload dto.ConsentPayload) error {
	if err := s.validateConsentPayload(payload); err != nil {
		return err
	}
	userID = strings.TrimSpace(userID)
	tel = strings.TrimSpace(tel)
	ip = strings.TrimSpace(ip)
	userAgent = strings.TrimSpace(userAgent)

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	var telPtr *string
	if tel != "" {
		telPtr = &tel
	}
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	var uaPtr *string
	if userAgent != "" {
		uaPtr = &userAgent
	}

	entries := []entity.ConsentLog{
		{
			UserID:        userIDPtr,
			Tel:           telPtr,
			ConsentType:   "privacy_policy",
			PolicyVersion: payload.PrivacyPolicyVersion,
			IPAddress:     ipPtr,
			UserAgent:     uaPtr,
		},
		{
			UserID:        userIDPtr,
			Tel:           telPtr,
			ConsentType:   "terms_of_service",
			PolicyVersion: payload.TermsVersion,
			IPAddress:     ipPtr,
			UserAgent:     uaPtr,
		},
	}
	for _, row := range entries {
		if err := s.repo.InsertConsentLog(ctx, row); err != nil {
			return err
		}
	}
	return nil
}

func (s privacyService) GetPolicyInfo() dto.PrivacyPolicyResponse {
	return dto.PrivacyPolicyResponse{
		PrivacyPolicyVersion: privacy.PrivacyPolicyVersion,
		TermsVersion:         privacy.TermsVersion,
		UpdatedAt:            privacy.PrivacyPolicyVersion,
		DPOEmail:             "privacy@pramool.in.th",
	}
}

func (s privacyService) CreateDSARRequest(ctx context.Context, userID string, req dto.CreateDSARRequest) (*dto.DSARRequestItem, error) {
	userID = strings.TrimSpace(userID)
	reqType := strings.TrimSpace(req.RequestType)
	switch reqType {
	case "access", "delete", "correct":
	default:
		return nil, ErrInvalidDSARType
	}
	note := strings.TrimSpace(req.Note)
	var notePtr *string
	if note != "" {
		notePtr = &note
	}
	row, err := s.repo.CreateDSARRequest(ctx, entity.DSARRequest{
		UserID:      userID,
		RequestType: reqType,
		Status:      "pending",
		UserNote:    notePtr,
	})
	if err != nil {
		return nil, err
	}
	if reqType == "access" {
		if err := s.generateAccessExport(ctx, userID, row.ID); err != nil {
			return nil, err
		}
		row.Status = "completed"
		now := time.Now()
		row.CompletedAt = &now
	}
	item := dsarItemFromEntity(*row)
	if ready, _ := s.repo.DSARExportExists(ctx, row.ID); ready {
		item.ExportReady = true
	}
	return item, nil
}

func (s privacyService) ListMyDSARRequests(ctx context.Context, userID string) ([]dto.DSARRequestItem, error) {
	rows, err := s.repo.ListDSARRequestsByUser(ctx, strings.TrimSpace(userID), 20)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DSARRequestItem, 0, len(rows))
	for _, row := range rows {
		item := *dsarItemFromEntity(row)
		if ready, _ := s.repo.DSARExportExists(ctx, row.ID); ready {
			item.ExportReady = true
		}
		items = append(items, item)
	}
	return items, nil
}

func (s privacyService) PurgeStaleTelVerify(ctx context.Context) (int64, error) {
	return s.repo.PurgeStaleTelVerify(ctx, privacy.TelVerifyRetention)
}

func dsarItemFromEntity(row entity.DSARRequest) *dto.DSARRequestItem {
	item := &dto.DSARRequestItem{
		ID:          row.ID,
		RequestType: row.RequestType,
		Status:      row.Status,
		CreatedAt:   row.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   row.UpdatedAt.Format(time.RFC3339),
	}
	if row.UserNote != nil {
		item.UserNote = *row.UserNote
	}
	if row.AdminNote != nil {
		item.AdminNote = *row.AdminNote
	}
	if row.CompletedAt != nil {
		item.CompletedAt = row.CompletedAt.Format(time.RFC3339)
	}
	if row.DueAt != nil {
		item.DueAt = row.DueAt.Format(time.RFC3339)
	}
	if row.DeletionExecutedAt != nil {
		item.DeletionExecutedAt = row.DeletionExecutedAt.Format(time.RFC3339)
	}
	return item
}
