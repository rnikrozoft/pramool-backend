package service

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

type BankService interface {
	ListActiveBanks(ctx context.Context) ([]entity.Bank, error)
}

type bankService struct {
	bankRepository repository.BankRepository
}

func NewBankService(bankRepository repository.BankRepository) BankService {
	return bankService{bankRepository: bankRepository}
}

func (s bankService) ListActiveBanks(ctx context.Context) ([]entity.Bank, error) {
	return s.bankRepository.ListActive(ctx)
}
