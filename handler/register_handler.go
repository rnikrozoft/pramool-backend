package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/mapping"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type RegisterHandler struct {
	validate               *validator.Validate
	authenticationService  service.AuthenticationService
	registerService        service.RegisterService
	accessCookieMaxAgeSec  int
	refreshCookieMaxAgeSec int
}

func NewRegisterHandler(
	validate *validator.Validate,
	authenticationService service.AuthenticationService,
	registerService service.RegisterService,
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
		accessCookieMaxAgeSec:  accessCookieMaxAgeSec,
		refreshCookieMaxAgeSec: refreshCookieMaxAgeSec,
	}
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

	userEntity := mapping.ToUserEntity(*user)
	if err := h.registerService.RegisterUser(ctx, userEntity); err != nil {
		return responseCommonError(c, err)
	}

	access, err := h.authenticationService.GenerateAccessToken(user.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	refresh, err := h.authenticationService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}
	ApplyAuthCookies(c, access, refresh, h.accessCookieMaxAgeSec, h.refreshCookieMaxAgeSec)
	return c.SendStatus(fiber.StatusCreated)
}
