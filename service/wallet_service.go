package service

type WalletService interface {
	WithdrawAllBalance(userID int) error
}
