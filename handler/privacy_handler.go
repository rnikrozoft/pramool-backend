package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/exception"
	"github.com/rnikrozoft/pramool-core/internal/privacy"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/repository"
	"github.com/rnikrozoft/pramool-core/service"
)

type PrivacyHandler struct {
	validate       *validator.Validate
	privacyService service.PrivacyService
}

func NewPrivacyHandler(validate *validator.Validate, privacyService service.PrivacyService) PrivacyHandler {
	return PrivacyHandler{
		validate:       validate,
		privacyService: privacyService,
	}
}

func (h PrivacyHandler) GetPolicyInfo(c *fiber.Ctx) error {
	return c.JSON(h.privacyService.GetPolicyInfo())
}

func (h PrivacyHandler) ListMyDSARRequests(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	items, err := h.privacyService.ListMyDSARRequests(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(dto.DSARRequestListResponse{Items: items})
}

func (h PrivacyHandler) CreateDSARRequest(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	req := new(dto.CreateDSARRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(exception.BadRequest(err))
	}
	item, err := h.privacyService.CreateDSARRequest(c.Context(), userID, *req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDSARType) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ประเภทคำขอไม่ถูกต้อง"})
		}
		return responseCommonError(c, err)
	}
	if req.RequestType == "access" {
		return c.Status(fiber.StatusCreated).JSON(item)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h PrivacyHandler) GetDeletionReadiness(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	resp, err := h.privacyService.GetAccountDeletionReadiness(c.Context(), userID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.JSON(resp)
}

func (h PrivacyHandler) ExecuteAccountDeletion(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	req := new(dto.ExecuteAccountDeletionRequest)
	_ = c.BodyParser(req)
	if err := h.privacyService.ExecuteAccountDeletion(c.Context(), userID, req.DSARRequestID); err != nil {
		if errors.Is(err, service.ErrAccountDeletionBlocked) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ยังไม่สามารถลบบัญชีได้ มีข้อผูกพันค้างอยู่"})
		}
		if errors.Is(err, service.ErrAccountAlreadyDeleted) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "บัญชีถูกลบแล้ว"})
		}
		return responseCommonError(c, err)
	}
	ClearAuthCookies(c)
	return c.JSON(dto.ExecuteAccountDeletionResponse{OK: true, UserID: userID, Message: "บัญชีถูกลบและทำให้ไม่ระบุตัวตนแล้ว"})
}

func (h PrivacyHandler) DownloadDSARExport(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	dsarID, err := c.ParamsInt("id")
	if err != nil || dsarID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid dsar id"})
	}
	raw, err := h.privacyService.GetDSARExport(c.Context(), userID, int64(dsarID))
	if err != nil {
		if errors.Is(err, repository.ErrPrivacyNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "ไม่พบไฟล์ส่งออก"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", "attachment; filename=\"pramool-data-export-"+strconv.Itoa(dsarID)+".json\"")
	return c.Send(raw)
}

func RequestMetaFromFiber(c *fiber.Ctx) (ip, userAgent string) {
	ip = strings.TrimSpace(c.IP())
	if forwarded := strings.TrimSpace(c.Get("X-Forwarded-For")); forwarded != "" {
		if idx := strings.Index(forwarded, ","); idx > 0 {
			ip = strings.TrimSpace(forwarded[:idx])
		} else {
			ip = forwarded
		}
	}
	userAgent = strings.TrimSpace(c.Get("User-Agent"))
	return ip, userAgent
}

func ValidateNationalID(userID string) error {
	if !privacy.IsValidThaiNationalID(userID) {
		return errors.New("เลขบัตรประชาชนไม่ถูกต้อง")
	}
	return nil
}
