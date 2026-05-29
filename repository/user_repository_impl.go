package repository

import (
	"database/sql"
	"errors"

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
	// Menggunakan LEFT JOIN agar pengguna tanpa profil (seperti Admin) tetap bisa dibaca secara aman.
	// COALESCE digunakan untuk memberikan nilai default (string kosong) apabila kolom profil bernilai NULL.
	query := `
		SELECT 
			u.id, u.email, u.role, u.status,
			COALESCE(p.full_name, '') AS full_name, 
			COALESCE(p.phone_number, '') AS phone_number, 
			COALESCE(p.nik, '') AS nik, 
			COALESCE(p.member_type, '') AS member_type, 
			COALESCE(p.address, '') AS address,
			COALESCE(p.photo_ktp_url, '') AS photo_ktp_url, 
			COALESCE(p.photo_selfie_url, '') AS photo_selfie_url, 
			COALESCE(p.bank_name, '') AS bank_name, 
			COALESCE(p.bank_account_number, '') AS bank_account_number,
			COALESCE((SELECT COUNT(*) FROM user_profiles WHERE referral_number = p.phone_number), 0) AS total_referral
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.id = ?`

	user := &model.UserResponse{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,                        // 1
		&user.Email,                     // 2
		&user.Role,                      // 3
		&user.Status,                    // 4
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

func (r *UserRepositoryImpl) FindUserByEmail(email string) (*model.UserResponse, error) {
	// Menggunakan LEFT JOIN agar pengguna tanpa profil (seperti Admin) tetap bisa dibaca secara aman.
	// COALESCE digunakan untuk memberikan nilai default (string kosong) apabila kolom profil bernilai NULL.
	query := `
		SELECT 
			u.id, u.email, u.role, u.status,
			COALESCE(p.full_name, '') AS full_name, 
			COALESCE(p.phone_number, '') AS phone_number, 
			COALESCE(p.nik, '') AS nik, 
			COALESCE(p.member_type, '') AS member_type, 
			COALESCE(p.address, '') AS address,
			COALESCE(p.photo_ktp_url, '') AS photo_ktp_url, 
			COALESCE(p.photo_selfie_url, '') AS photo_selfie_url, 
			COALESCE(p.bank_name, '') AS bank_name, 
			COALESCE(p.bank_account_number, '') AS bank_account_number,
			COALESCE((SELECT COUNT(*) FROM user_profiles WHERE referral_number = p.phone_number), 0) AS total_referral
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.email = ?`

	user := &model.UserResponse{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,                        // 1
		&user.Email,                     // 2
		&user.Role,                      // 3
		&user.Status,                    // 4
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

func (r *UserRepositoryImpl) GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error) {
	dashboard := &model.UserDashboardResponse{}

	// ==========================================
	// QUERY 1: Ambil Data Profil, Dompet, dan Total Statistik
	// ==========================================
	queryProfile := `
		SELECT 
			u.id, u.email, 
			p.full_name, p.phone_number,
			COALESCE(w.referral_balance, 0) AS current_balance,
			
			-- Subquery 1: Menghitung total orang yang daftar pakai nomor dia
			(SELECT COUNT(*) FROM user_profiles WHERE referral_number = p.phone_number) AS total_referred,
			
			-- Subquery 2: Menghitung total uang masuk yang berstatus success
			(SELECT COALESCE(SUM(amount), 0) FROM referral_rewards WHERE referrer_user_id = u.id AND status = 'success') AS total_earned
			
		FROM users u
		INNER JOIN user_profiles p ON u.id = p.user_id
		LEFT JOIN user_wallets w ON u.id = w.user_id
		WHERE u.id = ?`

	err := r.DB.QueryRow(queryProfile, userID).Scan(
		&dashboard.ID,
		&dashboard.Email,
		&dashboard.FullName,
		&dashboard.PhoneNumber,
		&dashboard.CurrentBalance,
		&dashboard.TotalReferredUsers,
		&dashboard.TotalEarnedReward,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	// ==========================================
	// QUERY 2: Ambil Daftar Orang yang Diajak (List Referrals)
	// ==========================================
	queryList := `
		SELECT u.id, p.full_name, p.phone_number, u.created_at
		FROM users u
		INNER JOIN user_profiles p ON u.id = p.user_id
		WHERE p.referral_number = ?
		ORDER BY u.created_at DESC`

	rowsList, err := r.DB.Query(queryList, dashboard.PhoneNumber)
	if err != nil {
		return nil, err
	}
	defer rowsList.Close()

	dashboard.ReferredUsersList = []model.ReferredUser{}

	for rowsList.Next() {
		var refUser model.ReferredUser
		err := rowsList.Scan(
			&refUser.ID,
			&refUser.FullName,
			&refUser.PhoneNumber,
			&refUser.RegisteredAt,
		)
		if err != nil {
			return nil, err
		}
		dashboard.ReferredUsersList = append(dashboard.ReferredUsersList, refUser)
	}

	// ==========================================
	// QUERY 3: Ambil Daftar Penarikan (Withdrawals)
	// ==========================================
	queryWithdrawals := `
		SELECT id, amount, status, created_at
		FROM referral_withdrawals
		WHERE user_id = ?
		ORDER BY created_at DESC`

	rowsWithdrawal, err := r.DB.Query(queryWithdrawals, userID)
	if err != nil {
		return nil, err
	}
	defer rowsWithdrawal.Close()

	dashboard.WithdrawalHistoryList = []model.WithdrawalHistory{}

	for rowsWithdrawal.Next() {
		var wHistory model.WithdrawalHistory
		err := rowsWithdrawal.Scan(
			&wHistory.ID,
			&wHistory.Amount,
			&wHistory.Status,
			&wHistory.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		dashboard.WithdrawalHistoryList = append(dashboard.WithdrawalHistoryList, wHistory)
	}

	return dashboard, nil
}
