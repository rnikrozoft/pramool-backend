package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/internal/privacy"
	"github.com/rnikrozoft/pramool-core/internal/retention"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

var (
	ErrAccountDeletionBlocked = errors.New("account deletion blocked")
	ErrAccountAlreadyDeleted  = errors.New("account already deleted")
	ErrInvalidCookieVersion   = errors.New("invalid cookie policy version")
)

func (s privacyService) RecordCookieConsent(ctx context.Context, userID, tel, ip, userAgent string, req dto.CookieConsentRequest) error {
	if req.CookiePolicyVersion != privacy.CookiePolicyVersion {
		return ErrInvalidCookieVersion
	}
	if !req.AcceptEssential {
		return errors.New("essential cookies required")
	}
	userID = strings.TrimSpace(userID)
	tel = strings.TrimSpace(tel)
	var userIDPtr, telPtr, ipPtr, uaPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if tel != "" {
		telPtr = &tel
	}
	if ip = strings.TrimSpace(ip); ip != "" {
		ipPtr = &ip
	}
	if userAgent = strings.TrimSpace(userAgent); userAgent != "" {
		uaPtr = &userAgent
	}
	if err := s.repo.InsertConsentLog(ctx, entity.ConsentLog{
		UserID: userIDPtr, Tel: telPtr, ConsentType: "essential_cookies",
		PolicyVersion: req.CookiePolicyVersion, IPAddress: ipPtr, UserAgent: uaPtr,
	}); err != nil {
		return err
	}
	if req.AcceptAnalytics {
		return s.repo.InsertConsentLog(ctx, entity.ConsentLog{
			UserID: userIDPtr, Tel: telPtr, ConsentType: "analytics_cookies",
			PolicyVersion: req.CookiePolicyVersion, IPAddress: ipPtr, UserAgent: uaPtr,
		})
	}
	return nil
}

func (s privacyService) GetMarketingOptIn(ctx context.Context, userID string) (bool, error) {
	return s.repo.GetMarketingOptIn(ctx, userID)
}

func (s privacyService) UpdateMarketingConsent(ctx context.Context, userID, tel, ip, userAgent string, req dto.MarketingConsentRequest) error {
	if req.PrivacyPolicyVersion != privacy.PrivacyPolicyVersion {
		return ErrInvalidPolicyVersion
	}
	if err := s.repo.SetMarketingOptIn(ctx, userID, req.MarketingOptIn); err != nil {
		return err
	}
	userID = strings.TrimSpace(userID)
	tel = strings.TrimSpace(tel)
	var userIDPtr, telPtr, ipPtr, uaPtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if tel != "" {
		telPtr = &tel
	}
	if ip = strings.TrimSpace(ip); ip != "" {
		ipPtr = &ip
	}
	if userAgent = strings.TrimSpace(userAgent); userAgent != "" {
		uaPtr = &userAgent
	}
	consentType := "marketing_opt_out"
	if req.MarketingOptIn {
		consentType = "marketing_opt_in"
	}
	return s.repo.InsertConsentLog(ctx, entity.ConsentLog{
		UserID: userIDPtr, Tel: telPtr, ConsentType: consentType,
		PolicyVersion: req.PrivacyPolicyVersion, IPAddress: ipPtr, UserAgent: uaPtr,
	})
}

func (s privacyService) GetAccountDeletionReadiness(ctx context.Context, userID string) (*dto.AccountDeletionBlockersResponse, error) {
	b, err := s.repo.GetAccountDeletionBlockers(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := &dto.AccountDeletionBlockersResponse{
		AlreadyDeleted:       b.AlreadyDeleted,
		CreditBalance:        b.CreditBalance,
		ActiveSellerAuctions: b.ActiveSellerAuctions,
		ActiveBidAuctions:    b.ActiveBidAuctions,
		PendingSellerShip:    b.PendingSellerShip,
		PendingBuyerConfirm:  b.PendingBuyerConfirm,
		PendingWithdrawals:   b.PendingWithdrawals,
	}
	var blockers []string
	if b.AlreadyDeleted {
		blockers = append(blockers, "บัญชีถูกลบแล้ว")
	}
	if b.CreditBalance > 0 {
		blockers = append(blockers, fmt.Sprintf("มียอดเครดิตคงเหลือ %d บาท", b.CreditBalance))
	}
	if b.ActiveSellerAuctions > 0 {
		blockers = append(blockers, fmt.Sprintf("มีประมูลที่เปิดอยู่ %d รายการ", b.ActiveSellerAuctions))
	}
	if b.ActiveBidAuctions > 0 {
		blockers = append(blockers, fmt.Sprintf("กำลังบิดอยู่ %d รายการ", b.ActiveBidAuctions))
	}
	if b.PendingSellerShip > 0 {
		blockers = append(blockers, fmt.Sprintf("รอจัดส่ง %d รายการ", b.PendingSellerShip))
	}
	if b.PendingBuyerConfirm > 0 {
		blockers = append(blockers, fmt.Sprintf("รอยืนยันรับของ %d รายการ", b.PendingBuyerConfirm))
	}
	if b.PendingWithdrawals > 0 {
		blockers = append(blockers, fmt.Sprintf("มีคำขอถอนเงินค้าง %d รายการ", b.PendingWithdrawals))
	}
	resp.Blockers = blockers
	resp.CanDelete = len(blockers) == 0
	return resp, nil
}

func (s privacyService) generateAccessExport(ctx context.Context, userID string, dsarID int64) error {
	payload, err := s.repo.BuildUserDataExport(ctx, userID)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := s.repo.SaveDSARExport(ctx, dsarID, raw); err != nil {
		return err
	}
	return s.repo.CompleteDSARRequest(ctx, dsarID)
}

func (s privacyService) GetDSARExport(ctx context.Context, userID string, dsarID int64) ([]byte, error) {
	req, err := s.repo.GetDSARRequestForUser(ctx, dsarID, userID)
	if err != nil {
		return nil, err
	}
	if req.RequestType != "access" {
		return nil, errors.New("export not available for this request type")
	}
	return s.repo.GetDSARExportJSON(ctx, dsarID, userID)
}

func (s privacyService) ExecuteAccountDeletion(ctx context.Context, userID string, dsarID int64) error {
	readiness, err := s.GetAccountDeletionReadiness(ctx, userID)
	if err != nil {
		return err
	}
	if readiness.AlreadyDeleted {
		return ErrAccountAlreadyDeleted
	}
	if !readiness.CanDelete {
		return ErrAccountDeletionBlocked
	}
	if err := s.repo.AnonymizeUserAccount(ctx, userID); err != nil {
		return err
	}
	if dsarID > 0 {
		_ = s.repo.MarkDSARDeletionExecuted(ctx, dsarID)
		_ = s.repo.CompleteDSARRequest(ctx, dsarID)
	}
	return nil
}

func (s privacyService) ListDataProcessors(ctx context.Context) ([]dto.DataProcessorItem, error) {
	rows, err := s.repo.ListActiveDataProcessors(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DataProcessorItem, 0, len(rows))
	for _, row := range rows {
		item := dto.DataProcessorItem{
			Name: row.Name, Purpose: row.Purpose, DataCategories: row.DataCategories,
			Location: row.Location, DPAStatus: row.DPAStatus,
		}
		if row.PrivacyURL != nil {
			item.PrivacyURL = *row.PrivacyURL
		}
		items = append(items, item)
	}
	return items, nil
}

func (s privacyService) ListRetentionJobs(ctx context.Context) (*dto.RetentionJobListResponse, error) {
	deps := &retention.Deps{Privacy: s.repo}
	inline := make([]dto.RetentionJobItem, 0)
	for _, job := range retention.InlineJobs(deps) {
		inline = append(inline, retentionJobDTO(job, true))
	}
	future := make([]dto.RetentionJobItem, 0)
	for _, job := range retention.FutureBatchJobs(deps) {
		future = append(future, retentionJobDTO(job, false))
	}
	dbRows, _ := s.repo.ListRetentionJobDefinitions(ctx)
	mergeRetentionRunMeta(inline, future, dbRows)
	return &dto.RetentionJobListResponse{
		InlineJobs: inline,
		FutureJobs: future,
		BatchNote:  "งานที่ batch_runner=future รอ batch runner กลางในอนาคต — ตอนนี้ยังไม่รันอัตโนมัติ",
	}, nil
}

func (s privacyService) RunInlineRetentionJob(ctx context.Context, jobID string) (int64, error) {
	deps := &retention.Deps{Privacy: s.repo}
	n, err := retention.RunInlineJob(ctx, jobID, deps)
	if err != nil {
		return 0, err
	}
	_ = s.repo.TouchRetentionJobRun(ctx, jobID, n)
	return n, nil
}

func (s privacyService) LogWinnerAddressDisclosure(ctx context.Context, sellerUserID, buyerUserID, auctionID, ip string) error {
	return s.repo.LogDataDisclosure(ctx, "winner_shipping_address", sellerUserID, buyerUserID, auctionID, ip)
}

func retentionJobDTO(job retention.Job, implemented bool) dto.RetentionJobItem {
	return dto.RetentionJobItem{
		JobID: job.ID, NameTH: job.NameTH, Description: job.Description,
		RetentionDays: job.RetentionDays, BatchRunner: job.BatchRunner,
		IsEnabled: job.Enabled, Implemented: implemented,
	}
}

func mergeRetentionRunMeta(inline, future []dto.RetentionJobItem, dbRows []repository.RetentionJobDefinitionRow) {
	byID := map[string]repository.RetentionJobDefinitionRow{}
	for _, row := range dbRows {
		byID[row.JobID] = row
	}
	apply := func(items []dto.RetentionJobItem) {
		for i := range items {
			if row, ok := byID[items[i].JobID]; ok {
				items[i].IsEnabled = row.IsEnabled
				if row.LastRunAt != nil {
					items[i].LastRunAt = row.LastRunAt.Format(time.RFC3339)
				}
				if row.LastDeletedCount != nil {
					items[i].LastDeletedCount = *row.LastDeletedCount
				}
			}
		}
	}
	apply(inline)
	apply(future)
}

func (s privacyService) IsAccountDeleted(ctx context.Context, userID string) (bool, error) {
	return s.repo.IsAccountDeleted(ctx, userID)
}

func (s privacyService) IsAccountDeletedByTel(ctx context.Context, tel string) (bool, error) {
	return s.repo.IsAccountDeletedByTel(ctx, tel)
}
