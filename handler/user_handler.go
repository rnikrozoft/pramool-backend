package handler

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/mapping"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type userHandler struct {
	validate    *validator.Validate
	userService service.UserService
}

func NewUserHandler(validate *validator.Validate, userService service.UserService) userHandler {
	return userHandler{
		validate:    validate,
		userService: userService,
	}
}

func (h userHandler) profileSubject(c *fiber.Ctx) (string, error) {
	sub, ok := c.Locals("user_id").(string)
	if !ok || sub == "" {
		return "", errors.New("cannot claims user information")
	}
	return sub, nil
}

func (h userHandler) GetMyInformation(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	u, err := h.userService.GetMyInfo(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	sn, bn, err := h.userService.CountUserFulfillmentBlocks(c.Context(), u.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	wb, wr := mapping.WithdrawalBlockedFromCounts(sn, bn)
	return c.JSON(mapping.ToUserProfileResponse(*u, wb, wr))
}

func (h userHandler) IsTelAlreadyUsed(c *fiber.Ctx) error {
	tel := c.Params("tel")
	if tel == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid tel"})
	}
	ok, err := h.userService.IsTelAlreadyUsed(c.Context(), tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.TelAvailabilityResponse{OK: ok})
}

func (h userHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}

	req := new(dto.UpdateProfileRequest)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	// JWT subject may still be tel (issued before registration). GET /users/profile resolves
	// tel → national user_id; UPDATE must use the same canonical id or no row matches.
	resolved, err := h.userService.GetMyInfo(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	canonicalID := strings.TrimSpace(resolved.UserID)
	if canonicalID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid user"})
	}

	if err := h.userService.UpdateMyProfile(c.Context(), mapping.UpdateProfileRequestToEntity(canonicalID, req)); err != nil {
		if errors.Is(err, service.ErrTelRequired) || errors.Is(err, service.ErrTelAlreadyUsed) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		if errors.Is(err, service.ErrUserNotRegistered) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "complete registration before updating profile"})
		}
		return responseCommonError(c, err)
	}

	u, err := h.userService.GetMyInfo(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	sn, bn, err := h.userService.CountUserFulfillmentBlocks(c.Context(), u.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	wb, wr := mapping.WithdrawalBlockedFromCounts(sn, bn)
	return c.JSON(mapping.ToUserProfileResponse(*u, wb, wr))
}

func (h userHandler) GetOnboardingStatus(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}

	isFirst, err := h.userService.IsFirstRegistration(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.OnboardingStatusResponse{IsFirstRegistration: isFirst})
}
