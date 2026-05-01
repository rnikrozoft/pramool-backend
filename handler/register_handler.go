package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/mapping"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type RegisterHandler struct {
	validate              *validator.Validate
	authenticationService service.AuthenticationService
	registerService       service.RegisterService
}

func NewRegisterHandler(
	validate *validator.Validate,
	authenticationService service.AuthenticationService,
	registerService service.RegisterService,
) RegisterHandler {
	return RegisterHandler{
		validate:              validate,
		authenticationService: authenticationService,
		registerService:       registerService,
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

	token, err := h.authenticationService.GenerateToken(user.UserID)
	if err != nil {
		return responseCommonError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   3600,
	})
	return c.SendStatus(fiber.StatusCreated)
}
