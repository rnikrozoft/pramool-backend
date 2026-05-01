package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type AuctionHandler struct {
	auctionService service.AuctionService
}

func NewAuctionHandler(auctionService service.AuctionService) AuctionHandler {
	return AuctionHandler{auctionService: auctionService}
}

func (h AuctionHandler) CreateAuction(c *fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(string)
	if !ok {
		return responseCommonError(c, errors.New("cannot claims user information"))
	}

	startPrice, err := strconv.ParseInt(c.FormValue("start_price"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid start_price"})
	}
	bidStep, err := strconv.ParseInt(c.FormValue("bid_step"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid bid_step"})
	}

	req := dto.CreateAuctionRequest{
		Title:       c.FormValue("title"),
		Category:    c.FormValue("category"),
		Condition:   c.FormValue("condition"),
		Description: c.FormValue("description"),
		StartPrice:  startPrice,
		BidStep:     bidStep,
		EndAt:       c.FormValue("end_at"),
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid multipart form"})
	}
	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "at least one image is required"})
	}
	if len(files) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "max 5 images"})
	}

	imagePaths := make([]string, 0, len(files))
	for i, file := range files {
		if err := validateImageFile(file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		fileDir := filepath.Join(".", "uploads", "auctions", sellerID)
		if err := os.MkdirAll(fileDir, 0o755); err != nil {
			return responseCommonError(c, err)
		}

		fileName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), i, sanitizeImageExt(file.Filename))
		savePath := filepath.Join(fileDir, fileName)
		if err := c.SaveFile(file, savePath); err != nil {
			return responseCommonError(c, err)
		}
		publicPath := "/uploads/auctions/" + sellerID + "/" + fileName
		imagePaths = append(imagePaths, publicPath)
	}

	result, err := h.auctionService.CreateAuction(c.Context(), sellerID, req, imagePaths)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h AuctionHandler) MyAuctions(c *fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(string)
	if !ok {
		return responseCommonError(c, errors.New("cannot claims user information"))
	}

	items, err := h.auctionService.GetSellerAuctions(c.Context(), sellerID)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": items})
}

func (h AuctionHandler) MyEarnings(c *fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(string)
	if !ok {
		return responseCommonError(c, errors.New("cannot claims user information"))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	items, err := h.auctionService.ListSellerEarnings(c.Context(), sellerID, limit, offset)
	if err != nil {
		return responseCommonError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": items, "limit": limit, "offset": offset})
}

func validateImageFile(file *multipart.FileHeader) error {
	maxBytes := int64(5 * 1024 * 1024)
	if file.Size > maxBytes {
		return fmt.Errorf("image %s exceeds 5MB", file.Filename)
	}

	contentType := strings.ToLower(file.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "image/jpeg") ||
		strings.HasPrefix(contentType, "image/jpg") ||
		strings.HasPrefix(contentType, "image/png") ||
		strings.HasPrefix(contentType, "image/webp") {
		return nil
	}
	return fmt.Errorf("unsupported image type for %s", file.Filename)
}

func sanitizeImageExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return ext
	default:
		return ".jpg"
	}
}
