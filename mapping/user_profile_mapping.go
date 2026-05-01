package mapping

import (
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
)

// ToUserProfileResponse maps a user row to the profile API shape.
func ToUserProfileResponse(u entity.User) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		UserID:         u.UserID,
		Tel:            u.Tel,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		AddressPrimary: u.AddressPrimary,
		Address:        u.Address,
		Soi:            u.Soi,
		Road:           u.Road,
		SubDistrict:    u.SubDistrict,
		District:       u.District,
		Province:       u.Province,
		ZipCode:        u.ZipCode,
		Email:          u.Email,
		Facebook:       u.Facebook,
		Credit:         u.Credit,
	}
}

// UpdateProfileRequestToEntity maps PUT body + JWT subject to a profile update row.
func UpdateProfileRequestToEntity(jwtUserID string, req *dto.UpdateProfileRequest) entity.ProfileUpdate {
	if req == nil {
		return entity.ProfileUpdate{UserID: jwtUserID}
	}
	return entity.ProfileUpdate{
		UserID:         jwtUserID,
		Tel:            req.Tel,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		AddressPrimary: req.AddressPrimary,
		Address:        req.Address,
		Soi:            req.Soi,
		Road:           req.Road,
		SubDistrict:    req.SubDistrict,
		District:       req.District,
		Province:       req.Province,
		ZipCode:        req.ZipCode,
		Email:          req.Email,
		Facebook:       req.Facebook,
	}
}
