package repository

type WalletRepository interface {
	WithdrawAllBalance(userID int) error
}
