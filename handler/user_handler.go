package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/constant"
	"github.com/rnikrozoft/pramool.in.th-backend/mapping"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{
		userService: userService,
	}
}

func (h userHandler) IsTelAlreadyUsed(c *fiber.Ctx) error {
	ok, err := h.userService.IsTelAlreadyUsed(c.Context(), c.Params("tel"))
	if err != nil {
		return responseCommonError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(map[string]bool{
		"ok": ok,
	})
}

func (h userHandler) GetMyInformation(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.ClaimsJWTObject).(string)
	if !ok {
		return responseCommonError(c, errors.New("cannot claims user information"))
	}

	userEntity, err := h.userService.GetMyInfo(c.Context(), claims)
	if err != nil {
		return responseCommonError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(mapping.ToUserDTO(*userEntity))
}
