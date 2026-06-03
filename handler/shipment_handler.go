package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type ShipmentHandler struct {
	svc service.ShipmentService
}

func NewShipmentHandler(svc service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{svc: svc}
}

func (h *ShipmentHandler) ListCarriers(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": h.svc.ListCarriers(),
	})
}

func (h *ShipmentHandler) GetWinnerShippingAddress(c *fiber.Ctx) error {
	auctionID := strings.TrimSpace(c.Params("id"))
	userID, _ := c.Locals("user_id").(string)
	if strings.TrimSpace(userID) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	if auctionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid auction"})
	}
	ip, _ := RequestMetaFromFiber(c)
	result, err := h.svc.GetWinnerShippingAddress(c.Context(), auctionID, userID, ip)
	if err != nil {
		if errors.Is(err, service.ErrMarkShippedNotAllowed) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ไม่สามารถดูที่อยู่ผู้ชนะได้"})
		}
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *ShipmentHandler) MarkSellerShipped(c *fiber.Ctx) error {
	auctionID := strings.TrimSpace(c.Params("id"))
	userID, _ := c.Locals("user_id").(string)
	if strings.TrimSpace(userID) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	if auctionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid auction"})
	}
	var body dto.MarkShippedRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}
	if err := h.svc.MarkSellerShipped(c.Context(), auctionID, userID, body.CarrierCode, body.TrackingNumber); err != nil {
		if errors.Is(err, service.ErrMarkShippedNotAllowed) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ไม่สามารถบันทึกการจัดส่งได้"})
		}
		if errors.Is(err, service.ErrInvalidShipmentCarrier) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "กรุณาเลือกขนส่ง"})
		}
		if errors.Is(err, service.ErrInvalidTrackingNumber) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "กรุณากรอกเลขพัสดุ"})
		}
		msg := err.Error()
		if msg != "" && !strings.Contains(strings.ToLower(msg), "trackingmore authentication") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg})
		}
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "บันทึกการจัดส่งแล้ว"})
}

func (h *ShipmentHandler) GetShipmentTracking(c *fiber.Ctx) error {
	auctionID := strings.TrimSpace(c.Params("id"))
	userID, _ := c.Locals("user_id").(string)
	if strings.TrimSpace(userID) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}
	if auctionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid auction"})
	}
	result, err := h.svc.RefreshTracking(c.Context(), auctionID, userID)
	if err != nil {
		if errors.Is(err, service.ErrShipmentAccessDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "ไม่มีสิทธิ์ดูข้อมูลพัสดุรายการนี้"})
		}
		if errors.Is(err, service.ErrShipmentNotRegistered) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ยังไม่มีข้อมูลการจัดส่ง"})
		}
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
		}
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
