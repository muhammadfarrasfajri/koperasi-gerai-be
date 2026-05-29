package service

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
)

type WalletServiceImpl struct {
	WalletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) *WalletServiceImpl {
	return &WalletServiceImpl{
		WalletRepo: walletRepo,
	}
}

func (s *WalletServiceImpl) WithdrawAllBalance(userID int) error {
	return s.WalletRepo.WithdrawAllBalance(userID)
}
