package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

// Parameter 'payment' diganti menjadi 'refresh' menyesuaikan isi kode di dalamnya
func (r *AuthRepositoryImpl) Register(user model.User, profile model.UserProfile, payment model.RegistrationPayment) error {

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Defer ini akan membatalkan semua perubahan jika terjadi error sebelum tx.Commit()
	defer tx.Rollback()

	// 1. Insert User
	queryUser := "INSERT INTO users (id, email, google_id, role, status) VALUES (?, ?, ?, ?, ?)"
	res, err := tx.Exec(queryUser, user.ID, user.Email, user.GoogleID, user.Role, user.Status)
	if err != nil {
		return err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// 2. CEK REFERRAL (Dipindah ke atas sebelum Insert Profile)
	var referredByID *int64 // Gunakan pointer agar bisa NULL di database jika tidak pakai referral

	if profile.ReferralNumber != "" {
		var tempID int64
		queryFindReferrer := "SELECT user_id FROM user_profiles WHERE phone_number = ?"
		err = tx.QueryRow(queryFindReferrer, profile.ReferralNumber).Scan(&tempID)

		if err == sql.ErrNoRows {
			// TOLAK registrasi jika nomor HP pengajak tidak terdaftar
			return errors.New("nomor HP referral tidak ditemukan atau tidak valid")
		} else if err != nil {
			return err // Error murni dari database
		}

		// Jika ketemu, simpan ID-nya
		referredByID = &tempID
	}

	// 3. Insert Profile (Ubah referral_number menjadi referred_by_id)
	queryProfile := `INSERT INTO user_profiles 
		(user_id, full_name, phone_number, nik, member_type, address, city, photo_ktp_url, photo_selfie_url, bank_name, bank_account_number, referred_by_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Kita masukkan referredByID (berisi angka/ID atau NULL) ke query
	_, err = tx.Exec(queryProfile, userID, profile.FullName, profile.PhoneNumber, profile.NIK, profile.MemberType, profile.Address, profile.City, profile.PhotoKTPURL, profile.PhotoSelfieURL, profile.BankName, profile.BankAccountNumber, referredByID)
	if err != nil {
		return err
	}

	// 4. Insert Payment Proof
	queryPayment := "INSERT INTO registration_payments (user_id, amount, payment_proof_url) VALUES (?, ?, ?)"
	_, err = tx.Exec(queryPayment, userID, payment.Amount, payment.PaymentProofURL)
	if err != nil {
		return err
	}

	// 5. Buat Dompet Kosong untuk Anggota Baru
	queryWallet := "INSERT INTO user_wallets (user_id, referral_balance, balance) VALUES (?, 0.00, 1000000.00)"
	_, err = tx.Exec(queryWallet, userID)
	if err != nil {
		return err
	}

	// 6. Eksekusi Komisi Instan Jika Pakai Referral (Lebih Ramping)
	if referredByID != nil {
		// A. Catat kas masuk dengan status langsung 'success'
		queryReward := "INSERT INTO referral_rewards (referrer_user_id, referred_user_id, amount, status) VALUES (?, ?, 100000.00, 'success')"
		_, errReward := tx.Exec(queryReward, *referredByID, userID)
		if errReward != nil {
			return errReward
		}

		// B. Langsung tambahkan saldo Rp 100.000 ke dompet pengajak
		queryUpdateWallet := "UPDATE user_wallets SET referral_balance = referral_balance + 100000.00 WHERE user_id = ?"
		_, errWallet := tx.Exec(queryUpdateWallet, *referredByID)
		if errWallet != nil {
			return errWallet
		}
	}

	// 7. Simpan Permanen (Semua query di atas disahkan secara bersamaan)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepositoryImpl) Login(rt model.RefreshToken) error {
	sqlQuery := `
        INSERT INTO refresh_tokens (user_id, token, expires_at) 
        VALUES (?, ?, ?) 
        ON DUPLICATE KEY UPDATE token = VALUES(token), expires_at = VALUES(expires_at)`
	_, err := r.DB.Exec(sqlQuery, rt.UserID, rt.Token, rt.ExpiresAt)
	return err
}

func (r *AuthRepositoryImpl) IsEmailExists(email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = ? LIMIT 1`

	var exists int
	err := r.DB.QueryRow(query, email).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *AuthRepositoryImpl) IsPhoneNoExists(phoneNumber string) (bool, error) {
	query := `SELECT 1 FROM user_profiles WHERE phone_number = ? LIMIT 1`

	var exists int
	err := r.DB.QueryRow(query, phoneNumber).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *AuthRepositoryImpl) IsNIKExists(nik string) (bool, error) {
	query := `SELECT 1 FROM user_profiles WHERE nik = ? LIMIT 1`

	var exists int
	err := r.DB.QueryRow(query, nik).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
