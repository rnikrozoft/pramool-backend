package mapping

import (
	"strings"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
)

// WithdrawalBlockedFromCounts builds withdrawal gating for GET /users (ถอนเงิน).
func WithdrawalBlockedFromCounts(sellerN, buyerN int) (blocked bool, reason string) {
	if sellerN == 0 && buyerN == 0 {
		return false, ""
	}
	var parts []string
	if sellerN > 0 {
		parts = append(parts, "กรุณาบันทึกการจัดส่งสินค้าให้ครบในฐานะผู้ขาย")
	}
	if buyerN > 0 {
		parts = append(parts, "กรุณายืนยันรับของให้ครบในฐานะผู้ชนะประมูล")
	}
	return true, strings.Join(parts, " และ ") + " ก่อนจึงจะถอนเงินได้"
}

// ToUserProfileResponse maps a user row to the profile API shape.
func ToUserProfileResponse(u entity.User, withdrawalBlocked bool, withdrawalReason string, pendingSellerShip int) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		UserID:                u.UserID,
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
		WithdrawalBlocked:        withdrawalBlocked,
		WithdrawalBlockReason:    withdrawalReason,
		PendingSellerShipCount: pendingSellerShip,
	}
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
