package repository

import (
	"context"

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
		query = query.Where("pay.status = ? OR u.is_active = ?", "pending", false)
	case "verified", "active":
		query = query.Where("u.is_active = ? AND pay.status = ?", true, "approved")
	}

	// 1. Hitung total baris yang cocok sebelum pagination
	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. Ambil data terpaginasi
	offset := (page - 1) * limit
	err = query.Select(`u.id as user_id, u.email, u.role, u.is_active, u.created_at,
			p.full_name, 
			p.phone_number as whatsapp_number, 
			p.phone_number as whats_app_number, 
			p.nik, p.member_type, p.address,
			p.photo_ktp_url, p.photo_selfie_url, p.bank_name, p.bank_account_number, p.referral_number,
			COALESCE(pay.amount, 0) as payment_amount, 
			COALESCE(pay.payment_proof_url, '') as payment_proof_url, 
			COALESCE(pay.status, '') as payment_status`).
		Limit(limit).
		Offset(offset).
		Order("u.created_at DESC").
		Scan(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (r *AdminRepositoryImpl) VerifyUser(ctx context.Context, userID int64, status string, isActive bool) error {
	return r.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update status pembayaran pendaftaran beserta waktu verifikasi (verified_at)
		updates := map[string]interface{}{
			"status":      status,
			"verified_at": gorm.Expr("NOW()"),
		}
		err := tx.Table("registration_payments").Where("user_id = ?", userID).Updates(updates).Error
		if err != nil {
			return err
		}

		// 2. Update status keaktifan user
		err = tx.Table("users").Where("id = ?", userID).Update("is_active", isActive).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *AdminRepositoryImpl) GetAllProfiles(ctx context.Context) ([]model.UserProfile, error) {
	var profiles []model.UserProfile
	
	err := r.GormDB.WithContext(ctx).Table("user_profiles").
		Select(`id, user_id, full_name, phone_number, nik, member_type, address, 
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
		Joins("LEFT JOIN registration_payments pay ON u.id = pay.user_id").
		Where("u.role = ? AND (pay.status = ? OR u.is_active = ?)", "member", "pending", false).
		Count(&summary.TotalPendingVerification).Error
	if err != nil {
		return summary, err
	}

	return summary, nil
}
