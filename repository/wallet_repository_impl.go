package repository

import (
	"database/sql"
	"errors"
)

type WalletRepositoryImpl struct {
	DB *sql.DB
}

func NewWalletRepositoryImpl(db *sql.DB) *WalletRepositoryImpl {
	return &WalletRepositoryImpl{
		DB: db,
	}
}

func (r *WalletRepositoryImpl) WithdrawAllBalance(userID int) error {
	// 1. Mulai Transaksi Database
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Otomatis batalkan semua perubahan jika terjadi error sebelum Commit
	defer tx.Rollback()

	// 2. Cek Saldo Saat Ini sekaligus Kunci Baris (Row-Level Lock)
	var currentBalance float64
	queryCheck := "SELECT referral_balance FROM user_wallets WHERE user_id = ? FOR UPDATE"

	err = tx.QueryRow(queryCheck, userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("dompet tidak ditemukan")
		}
		return err
	}

	// 3. Validasi: Tolak jika saldo kosong
	if currentBalance <= 0 {
		return errors.New("saldo kosong, tidak ada yang bisa ditarik")
	}

	// 4. Catat Riwayat Penarikan Langsung Sebesar Total Saldo
	queryWithdraw := "INSERT INTO referral_withdrawals (user_id, amount, status) VALUES (?, ?, 'pending')"
	_, err = tx.Exec(queryWithdraw, userID, currentBalance)
	if err != nil {
		return err
	}

	// 5. Kuras / Nol-kan Saldo di Dompet
	queryResetWallet := "UPDATE user_wallets SET referral_balance = 0 WHERE user_id = ?"
	_, err = tx.Exec(queryResetWallet, userID)
	if err != nil {
		return err
	}

	// 6. Simpan Permanen (Commit Transaksi & Lepas Kunci)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
