package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/service"
)

type bankHandler struct {
	bankService service.BankService
}

func NewBankHandler(bankService service.BankService) bankHandler {
	return bankHandler{bankService: bankService}
}

func (h bankHandler) List(c *fiber.Ctx) error {
	banks, err := h.bankService.ListActiveBanks(c.Context())
	if err != nil {
		return responseCommonError(c, err)
	}

	resp := make([]dto.BankOptionResponse, 0, len(banks))
	for _, bank := range banks {
		resp = append(resp, dto.BankOptionResponse{
			BankID:   bank.BankID,
			BankCode: bank.BankCode,
			NameTH:   bank.BankNameTH,
			NameEN:   bank.BankNameEN,
		})
	}
	return c.JSON(resp)
}
