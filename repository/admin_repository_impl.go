package repository

import (
	"context"
	"errors"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/config"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	GormDB *gorm.DB
}

// NewAdminRepository membungkus koneksi basis data sql.DB tim secara dinamis ke GORM
func NewAdminRepository() *AdminRepositoryImpl {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: config.DB, // Membungkus *sql.DB global dari config.DB
	}), &gorm.Config{})
	if err != nil {
		panic("failed to initialize GORM for AdminRepository: " + err.Error())
	}

	return &AdminRepositoryImpl{
		GormDB: gormDB,
	}
}

func (r *AdminRepositoryImpl) GetUsers(ctx context.Context, statusFilter string, page int, limit int) ([]model.UserDetail, int64, error) {
	var users []model.UserDetail
	var totalRows int64

	query := r.GormDB.WithContext(ctx).Table("users u").
		Joins("LEFT JOIN user_profiles p ON u.id = p.user_id").
		Joins("LEFT JOIN registration_payments pay ON u.id = pay.user_id").
		Where("u.role = ?", "member")

	switch statusFilter {
	case "pending":
		query = query.Where("u.status = ?", "pending")
	case "verified", "active", "approve":
		query = query.Where("u.status = ?", "active")
	case "reject", "rejected":
		query = query.Where("u.status = ?", "reject")
	}

	// 1. Hitung total baris yang cocok sebelum pagination
	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. Ambil data terpaginasi
	offset := (page - 1) * limit
	err = query.Select(`u.id as user_id, u.email, u.role, u.status, u.created_at,
			p.full_name, 
			p.phone_number as whatsapp_number, 
			p.nik, p.member_type, p.address, COALESCE(p.city, '') as city,
			p.photo_ktp_url, p.photo_selfie_url, p.bank_name, p.bank_account_number, p.referral_number,
			COALESCE(pay.amount, 0) as payment_amount, 
			COALESCE(pay.payment_proof_url, '') as payment_proof_url, 
			COALESCE(u.rejection_reason, '') as rejection_reason`).
		Limit(limit).
		Offset(offset).
		Order("u.created_at DESC").
		Scan(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (r *AdminRepositoryImpl) VerifyUser(ctx context.Context, userID int64, paymentStatus string, userStatus string, rejectionReason string, verifiedBy int64) error {
	return r.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cari admins.id berdasarkan users.id (verifiedBy)
		var adminID int64
		err := tx.Table("admins").
			Select("id").
			Where("email = (SELECT email FROM users WHERE id = ?)", verifiedBy).
			Scan(&adminID).Error
		if err != nil {
			return err
		}

		// Update kolom verifikasi (status, verified_at, verified_by, rejection_reason) langsung di tabel users
		updates := map[string]interface{}{
			"status":           userStatus,
			"verified_at":      gorm.Expr("NOW()"),
			"rejection_reason": rejectionReason,
		}
		if adminID > 0 {
			updates["verified_by"] = adminID
		}
		err = tx.Table("users").Where("id = ?", userID).Updates(updates).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *AdminRepositoryImpl) GetAllProfiles(ctx context.Context) ([]model.UserProfile, error) {
	var profiles []model.UserProfile

	err := r.GormDB.WithContext(ctx).Table("user_profiles").
		Select(`id, user_id, full_name, phone_number, nik, member_type, address, city, 
				photo_ktp_url, photo_selfie_url, bank_name, bank_account_number, 
				referral_number`).
		Scan(&profiles).Error

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *AdminRepositoryImpl) GetAdminSummary(ctx context.Context) (model.AdminSummary, error) {
	var summary model.AdminSummary

	// 1. Hitung total user dengan role 'member'
	err := r.GormDB.WithContext(ctx).Table("users").
		Where("role = ?", "member").
		Count(&summary.TotalUsers).Error
	if err != nil {
		return summary, err
	}

	// 2. Hitung total pengajuan pending verifikasi
	// (yaitu: users yang is_active = false atau registration_payments yang statusnya 'pending')
	err = r.GormDB.WithContext(ctx).Table("users u").
		Where("u.role = ? AND u.status = ?", "member", "pending").
		Count(&summary.TotalPendingVerification).Error
	if err != nil {
		return summary, err
	}

	return summary, nil
}

func (r *AdminRepositoryImpl) GetUserByID(ctx context.Context, userID int64) (model.UserDetail, error) {
	var user model.UserDetail

	err := r.GormDB.WithContext(ctx).Table("users u").
		Joins("LEFT JOIN user_profiles p ON u.id = p.user_id").
		Joins("LEFT JOIN registration_payments pay ON u.id = pay.user_id").
		Where("u.id = ? AND u.role = ?", userID, "member").
		Select(`u.id as user_id, u.email, u.role, u.status, u.created_at,
				p.full_name, 
				p.phone_number as whatsapp_number, 
				p.nik, p.member_type, p.address, COALESCE(p.city, '') as city,
				p.photo_ktp_url, p.photo_selfie_url, p.bank_name, p.bank_account_number, p.referral_number,
				COALESCE(pay.amount, 0) as payment_amount, 
				COALESCE(pay.payment_proof_url, '') as payment_proof_url, 
				COALESCE(u.rejection_reason, '') as rejection_reason`).
		Scan(&user).Error

	if err != nil {
		return user, err
	}
	if user.UserID == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}

// adminID didapatkan dari token JWT di Controller (ctx.GetInt("user_id"))
func (r *AdminRepositoryImpl) VerifyWithdrawal(adminID int, req model.VerifyWithdrawalRequest) error {
	// Membuka transaksi database ala GORM
	return r.GormDB.Transaction(func(tx *gorm.DB) error {

		// 1. Kunci data penarikan (Pessimistic Locking) dan cek status saat ini
		var withdrawal struct {
			UserID int
			Amount int
			Status string
		}

		// Gunakan tx.Raw().Scan() untuk raw query SELECT di GORM
		err := tx.Raw("SELECT user_id, amount, status FROM referral_withdrawals WHERE id = ? FOR UPDATE", req.ID).Scan(&withdrawal).Error
		if err != nil {
			return errors.New("gagal menemukan data penarikan")
		}

		if withdrawal.Status != "pending" {
			return errors.New("penarikan ini sudah diproses sebelumnya")
		}

		// 2. Query UPDATE Status Penarikan
		queryUpdate := `
			UPDATE referral_withdrawals 
			SET status = ?, reject_reason = ?, verified_at = NOW(), verified_by = ? 
			WHERE id = ?`

		// FIX GORM: Tangkap error dengan .Error di akhir fungsi
		if err := tx.Exec(queryUpdate, req.Status, req.RejectReason, adminID, req.ID).Error; err != nil {
			return err
		}

		// 3. Logika Refund: Kembalikan saldo jika penarikan ditolak
		if req.Status == "reject" {
			// Sesuaikan 'user_wallets' dengan tabel tempat kamu menyimpan saldo
			queryRefund := `UPDATE user_wallets SET referral_balance = referral_balance + ? WHERE user_id = ?`
			if err := tx.Exec(queryRefund, withdrawal.Amount, withdrawal.UserID).Error; err != nil {
				return errors.New("gagal mengembalikan saldo referral: " + err.Error())
			}
		}

		// Jika return nil, GORM akan otomatis melakukan tx.Commit()
		// Jika return error, GORM akan otomatis melakukan tx.Rollback()
		return nil
	})
}
