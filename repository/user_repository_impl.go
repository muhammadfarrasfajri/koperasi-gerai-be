package repository

import (
	"database/sql"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) FindUserById(id int) (*model.UserResponse, error) {
	// Ganti LEFT JOIN menjadi INNER JOIN untuk memastikan tidak ada data NULL yang merusak .Scan() Golang
	query := `
		SELECT 
			u.id, u.email, u.role, u.is_active,
			p.full_name, p.phone_number, p.nik, p.member_type, p.address,
			p.photo_ktp_url, p.photo_selfie_url, p.bank_name, p.bank_account_number,
			(SELECT COUNT(*) FROM user_profiles WHERE referral_number = p.phone_number) AS total_referral
		FROM users u
		INNER JOIN user_profiles p ON u.id = p.user_id
		WHERE u.id = ?`

	// Jika Profile di struct adalah pointer, gunakan inisialisasi ini:
	// user := &model.UserResponse{ Profile: &model.UserProfile{} }

	// Jika Profile BUKAN pointer, cukup seperti ini:
	user := &model.UserResponse{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.IsActive,
		&user.TotalReferral,
		&user.Profile.FullName,
		&user.Profile.PhoneNumber,
		&user.Profile.NIK,
		&user.Profile.MemberType,
		&user.Profile.Address,
		&user.Profile.PhotoKTPURL,
		&user.Profile.PhotoSelfieURL,
		&user.Profile.BankName,
		&user.Profile.BankAccountNumber,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindUserByEmail(email string) (*model.UserResponse, error) {
	query := `
        SELECT 
            u.id, u.email, u.role, u.is_active,
            p.full_name, p.phone_number, p.nik, p.member_type, p.address,
            p.photo_ktp_url, p.photo_selfie_url, p.bank_name, p.bank_account_number,
            (SELECT COUNT(*) FROM user_profiles WHERE referral_number = p.phone_number) AS total_referral
        FROM users u
        INNER JOIN user_profiles p ON u.id = p.user_id
        WHERE u.email = ?`

	user := &model.UserResponse{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,                        // 1
		&user.Email,                     // 2
		&user.Role,                      // 3
		&user.IsActive,                  // 4
		&user.Profile.FullName,          // 5
		&user.Profile.PhoneNumber,       // 6
		&user.Profile.NIK,               // 7
		&user.Profile.MemberType,        // 8
		&user.Profile.Address,           // 9
		&user.Profile.PhotoKTPURL,       // 10
		&user.Profile.PhotoSelfieURL,    // 11
		&user.Profile.BankName,          // 12
		&user.Profile.BankAccountNumber, // 13
		&user.TotalReferral,             // 14
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
