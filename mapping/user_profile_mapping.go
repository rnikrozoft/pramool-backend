package mapping

import (
	"math"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
)

func creditDebtBaht(credit int64) int64 {
	if credit < 0 {
		return -credit
	}
	return 0
}

// WithdrawalBlockedFromCounts builds withdrawal gating for GET /users (ถอนเงิน).
func WithdrawalBlockedFromCounts(pendingSellerShip int) (blocked bool, reason string) {
	if pendingSellerShip == 0 {
		return false, ""
	}
	return true, "กรุณาบันทึกการจัดส่งสินค้าให้ครบในฐานะผู้ขาย ก่อนจึงจะถอนเงินได้"
}

func sellerReviewAvgRating(starPoints int64, reviewCount int) float64 {
	if reviewCount <= 0 {
		return 0
	}
	avg := float64(starPoints) / float64(reviewCount) / 2.0
	if avg < 0 {
		avg = 0
	}
	if avg > 5 {
		avg = 5
	}
	return math.Round(avg*10) / 10
}

// ToUserProfileResponse maps a user row to the profile API shape.
func ToUserProfileResponse(u entity.User, withdrawalBlocked bool, withdrawalReason string, pendingSellerShip int, appealPending bool, appealStatus string, unreadNotifications int) dto.UserProfileResponse {
	restricted := u.RestrictedUntil != nil && u.RestrictedUntil.After(time.Now())
	postingRestricted := u.PostingRestrictedUntil != nil && u.PostingRestrictedUntil.After(time.Now())
	resp := dto.UserProfileResponse{
		UserID:                u.UserID,
		NationalID:            u.NationalID,
		Tel:                   u.Tel,
		FirstName:             u.FirstName,
		LastName:              u.LastName,
		AddressPrimary:        u.AddressPrimary,
		Address:               u.Address,
		Soi:                   u.Soi,
		Road:                  u.Road,
		SubDistrict:           u.SubDistrict,
		District:              u.District,
		Province:              u.Province,
		ZipCode:               u.ZipCode,
		Email:                 u.Email,
		Facebook:              u.Facebook,
		BankID:                u.BankID,
		BankAccountName:       u.BankAccountName,
		BankAccountNumber:     u.BankAccountNumber,
		Credit:                u.Credit,
		HasCreditDebt:         u.Credit < 0,
		CreditDebtBaht:        creditDebtBaht(u.Credit),
		WithdrawalBlocked:        withdrawalBlocked,
		WithdrawalBlockReason:    withdrawalReason,
		PendingSellerShipCount: pendingSellerShip,
		AccountRestricted:       restricted,
		PostingRestricted:       postingRestricted,
		ReputationPoints:        u.ReputationPoints,
		SellerReviewAvgRating:   sellerReviewAvgRating(u.ReputationPoints, u.SellerReviewCount),
		SellerReviewCount:       u.SellerReviewCount,
		AppealPending:           appealPending,
		AppealStatus:            appealStatus,
		UnreadNotificationCount: unreadNotifications,
	}
	if restricted && u.RestrictedUntil != nil {
		resp.RestrictedUntil = u.RestrictedUntil.Format(time.RFC3339)
	}
	if strings.TrimSpace(u.RestrictedReason) != "" {
		resp.RestrictedReason = strings.TrimSpace(u.RestrictedReason)
	}
	if postingRestricted && u.PostingRestrictedUntil != nil {
		resp.PostingRestrictedUntil = u.PostingRestrictedUntil.Format(time.RFC3339)
	}
	if strings.TrimSpace(u.PostingRestrictedReason) != "" {
		resp.PostingRestrictedReason = strings.TrimSpace(u.PostingRestrictedReason)
	}
	return resp
}

// UpdateProfileRequestToEntity maps PUT body + JWT subject to a profile update row.
func UpdateProfileRequestToEntity(jwtUserID string, req *dto.UpdateProfileRequest) entity.ProfileUpdate {
	if req == nil {
		return entity.ProfileUpdate{UserID: jwtUserID}
	}
	return entity.ProfileUpdate{
		UserID:            jwtUserID,
		Tel:               req.Tel,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		AddressPrimary:    req.AddressPrimary,
		Address:           req.Address,
		Soi:               req.Soi,
		Road:              req.Road,
		SubDistrict:       req.SubDistrict,
		District:          req.District,
		Province:          req.Province,
		ZipCode:           req.ZipCode,
		Email:             req.Email,
		Facebook:          req.Facebook,
		BankID:            req.BankID,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
	}
}
