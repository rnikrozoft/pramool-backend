package handler

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/exception"
	"github.com/rnikrozoft/pramool-core/mapping"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type RegisterHandler struct {
	validate               *validator.Validate
	authenticationService  service.AuthenticationService
	registerService        service.RegisterService
	userService            service.UserService
	privacyService         service.PrivacyService
	accessCookieMaxAgeSec  int
	refreshCookieMaxAgeSec int
}

func NewRegisterHandler(
	validate *validator.Validate,
	authenticationService service.AuthenticationService,
	registerService service.RegisterService,
	userService service.UserService,
	privacyService service.PrivacyService,
	accessCookieMaxAgeSec int,
	refreshCookieMaxAgeSec int,
) RegisterHandler {
	if accessCookieMaxAgeSec <= 0 {
		accessCookieMaxAgeSec = 3600
	}
	if refreshCookieMaxAgeSec <= 0 {
		refreshCookieMaxAgeSec = 3600 * 24 * 7
	}
	return RegisterHandler{
		validate:               validate,
		authenticationService:  authenticationService,
		registerService:        registerService,
		userService:            userService,
		privacyService:         privacyService,
		accessCookieMaxAgeSec:  accessCookieMaxAgeSec,
		refreshCookieMaxAgeSec: refreshCookieMaxAgeSec,
	}
}

// Signup stores a bcrypt password on tel_verify and logs the user in (sets auth cookies).
func (h RegisterHandler) Signup(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(dto.SignupRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Tel = strings.TrimSpace(req.Tel)
	req.Email = strings.TrimSpace(req.Email)
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	used, err := h.userService.IsTelAlreadyUsed(ctx, req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	if used {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "เบอร์โทรศัพท์นี้ลงทะเบียนแล้ว"})
	}
	hasPW, err := h.registerService.TelVerifyHasPassword(ctx, req.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	if hasPW {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "เบอร์โทรศัพท์นี้มีบัญชีแล้ว กรุณาเข้าสู่ระบบเพื่อทำรายการต่อ"})
	}
	if err := h.registerService.RegisterTelIfNotExist(ctx, req.Tel); err != nil {
		return responseCommonError(c, err)
	}
	if err := h.registerService.SignupWithPassword(ctx, req.FirstName, req.LastName, req.Tel, req.Email, req.Password); err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
		}
		return responseCommonError(c, err)
	}
	ip, ua := RequestMetaFromFiber(c)
	if err := h.privacyService.RecordConsents(ctx, "", req.Tel, ip, ua, req.ConsentPayload); err != nil {
		if errors.Is(err, service.ErrInvalidPolicyVersion) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "เวอร์ชันนโยบายความเป็นส่วนตัวไม่ตรงกับปัจจุบัน กรุณารีเฟรชหน้าแล้วลองใหม่"})
		}
		return responseCommonError(c, err)
	}
	tokens, err := h.authenticationService.LoginByTel(ctx, req.Tel, req.Password, false)
	if err != nil {
		return responseCommonError(c, err)
	}
	accessSec, refreshSec := SessionCookieMaxAges(false, h.accessCookieMaxAgeSec, h.refreshCookieMaxAgeSec)
	ApplyAuthCookies(c, tokens.Access, tokens.Refresh, accessSec, refreshSec)
	return c.SendStatus(fiber.StatusCreated)
}

// Register godoc
// @Summary      Register new user and return authentication token
// @Description  Create a new user and return JWT token upon success
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserRegisterRequest  true  "User data"
// @Success      201   {object}  model.RegisterResponse
// @Failure      400   {object}  exception.Response "Invalid input or registration failure"
// @Failure      500   {object}  exception.Response "Cannot register or token generation failure"
// @Router       /register [post]
func (h RegisterHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()

	user := new(dto.UserRegisterRequest)
	if err := validate(c, h.validate, user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	user.NationalID = strings.TrimSpace(user.NationalID)
	if user.NationalID == "" {
		user.NationalID = strings.TrimSpace(user.UserID)
	}
	user.Tel = strings.TrimSpace(user.Tel)
	user.Email = strings.TrimSpace(user.Email)
	if user.NationalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "กรุณากรอกเลขบัตรประชาชน"})
	}
	if err := ValidateNationalID(user.NationalID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	userEntity := mapping.ToUserEntity(*user)
	if err := h.registerService.RegisterUser(ctx, userEntity); err != nil {
		if errors.Is(err, service.ErrNationalIDAlreadyRegistered) ||
			errors.Is(err, service.ErrEmailAlreadyRegistered) ||
			errors.Is(err, service.ErrTelHasFullUserRecord) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
		}
		return responseCommonError(c, err)
	}
	internalUserID, err := h.userService.FindUserIDByTelWithFallback(ctx, user.Tel)
	if err != nil {
		return responseCommonError(c, err)
	}
	ip, ua := RequestMetaFromFiber(c)
	if err := h.privacyService.RecordConsents(ctx, internalUserID, user.Tel, ip, ua, user.ConsentPayload); err != nil {
		if errors.Is(err, service.ErrInvalidPolicyVersion) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "เวอร์ชันนโยบายความเป็นส่วนตัวไม่ตรงกับปัจจุบัน กรุณารีเฟรชหน้าแล้วลองใหม่"})
		}
		return responseCommonError(c, err)
	}

	access, err := h.authenticationService.GenerateAccessToken(internalUserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	refresh, err := h.authenticationService.GenerateRefreshTokenWithRemember(internalUserID, false)
	if err != nil {
		return responseCommonError(c, err)
	}
	accessSec, refreshSec := SessionCookieMaxAges(false, h.accessCookieMaxAgeSec, h.refreshCookieMaxAgeSec)
	ApplyAuthCookies(c, access, refresh, accessSec, refreshSec)
	return c.SendStatus(fiber.StatusCreated)
}
