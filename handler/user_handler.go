package handler

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/mapping"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
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

func (h userHandler) profileResponse(c *fiber.Ctx, u *entity.User) error {
	sn, err := h.userService.CountUserFulfillmentBlocks(c.Context(), u.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	wb, wr := mapping.WithdrawalBlockedFromCounts(sn)
	appealPending, appealStatus, err := h.userService.ProfileAppealMeta(c.Context(), u.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	unread, err := h.userService.CountUnreadNotifications(c.Context(), u.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(mapping.ToUserProfileResponse(*u, wb, wr, sn, appealPending, appealStatus, unread))
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
	return h.profileResponse(c, u)
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
	return h.profileResponse(c, u)
}

func (h userHandler) SubmitRestrictionAppeal(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	req := new(dto.SubmitRestrictionAppealRequest)
	if err := validate(c, h.validate, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	appealID, err := h.userService.SubmitRestrictionAppeal(c.Context(), userID, req.Reason)
	if err != nil {
		if errors.Is(err, repository.ErrAppealNotRestricted) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "บัญชีไม่ได้ถูกจำกัดอยู่"})
		}
		if errors.Is(err, repository.ErrAppealDuplicate) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "มีคำขอที่รอตรวจสอบอยู่แล้ว"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.RestrictionAppealResponse{
		AppealID: appealID,
		Status:   "pending",
	})
}

func (h userHandler) GetRestrictionAppeal(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	row, err := h.userService.GetRestrictionAppeal(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	if row == nil {
		return c.JSON(dto.RestrictionAppealResponse{Status: "none"})
	}
	resp := dto.RestrictionAppealResponse{
		AppealID: row.AppealID,
		Status:   row.Status,
		Reason:   row.Reason,
		CreatedAt: row.CreatedAt.Format(time.RFC3339),
		AdminNote: repository.AppealNoteString(row.AdminNote),
	}
	if row.ResolvedAt != nil {
		resp.ResolvedAt = row.ResolvedAt.Format(time.RFC3339)
	}
	return c.JSON(resp)
}

func (h userHandler) GetOnboardingStatus(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}

	status, err := h.userService.GetOnboardingStatus(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(status)
}

func (h userHandler) ListNotifications(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	rows, total, err := h.userService.ListNotifications(c.Context(), userID, limit, offset)
	if err != nil {
		return responseCommonError(c, err)
	}
	items := make([]dto.UserNotificationItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapping.ToUserNotificationItem(row))
	}
	return c.JSON(dto.UserNotificationListResponse{Items: items, Total: total})
}

func (h userHandler) GetUnreadNotificationCount(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	count, err := h.userService.CountUnreadNotifications(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.UnreadNotificationCountResponse{Count: count})
}

func (h userHandler) MarkNotificationRead(c *fiber.Ctx) error {
	userID, err := h.profileSubject(c)
	if err != nil {
		return responseCommonError(c, err)
	}
	notificationID, err := c.ParamsInt("id")
	if err != nil || notificationID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid notification id"})
	}
	row, err := h.userService.MarkNotificationRead(c.Context(), userID, int64(notificationID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
		}
		return responseCommonError(c, err)
	}
	readAt := ""
	if row.ReadAt != nil {
		readAt = row.ReadAt.Format(time.RFC3339)
	}
	return c.JSON(dto.MarkNotificationReadResponse{
		ReadAt:         readAt,
		ExpiresAt:      row.ExpiresAt.Format(time.RFC3339),
		AutoDeleteNote:   dto.NotificationAutoDeleteNote,
	})
}
