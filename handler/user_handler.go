package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool.in.th-backend/constant"
	"github.com/rnikrozoft/pramool.in.th-backend/exception"
	"github.com/rnikrozoft/pramool.in.th-backend/mapping"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (h UserHandler) GetMyInformation(c *fiber.Ctx) error {
	claims, ok := c.Locals("userClaims").(*model.CustomClaims)
	if !ok {
		e := exception.Set(constant.ErrorSomethingWentWrong)
		return c.Status(fiber.StatusUnauthorized).JSON(e)
	}

	userEntity, err := h.userService.GetMyInfo(c.Context(), claims.UserID)
	if err != nil {
		e := exception.Set(constant.ErrorSomethingWentWrong)
		return c.Status(fiber.StatusInternalServerError).JSON(e)
	}

	return c.Status(fiber.StatusOK).JSON(mapping.ToUserDTO(*userEntity))
}
