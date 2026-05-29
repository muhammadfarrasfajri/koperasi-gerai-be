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

func (r *UserRepositoryImpl) FindUserById(id int, role string) (*model.UserResponse, error) {
	// ==========================================
	// 1. JIKA ROLE ADALAH ADMIN
	// ==========================================
	if role == "admin" {
		var email string
		err := r.DB.QueryRow("SELECT email FROM users WHERE id = ?", id).Scan(&email)
		if err != nil {
			return nil, err
		}

		queryAdmin := `
			SELECT id, email, full_name, role, is_active
			FROM admins
			WHERE email = ?
			LIMIT 1`
		var admin struct {
			ID       int
			Email    string
			FullName string
			Role     string
			IsActive bool
		}
		err = r.DB.QueryRow(queryAdmin, email).Scan(&admin.ID, &admin.Email, &admin.FullName, &admin.Role, &admin.IsActive)
		if err != nil {
			return nil, err
		}

		status := "active"
		if !admin.IsActive {
			status = "pending"
		}

		return &model.UserResponse{
			ID:     id,
			Email:  admin.Email,
			Role:   admin.Role,
			Status: status,
			Profile: model.UserProfile{
				FullName: admin.FullName,
			},
		}, nil
	}

	// ==========================================
	// 2. JIKA ROLE BUKAN ADMIN (USER BIASA)
	// ==========================================
	query := `
		SELECT 
			u.id, u.email, u.role, u.status, u.rejection_reason,
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
		&user.RejectionReason,           // 5
		&user.Profile.FullName,          // 6
		&user.Profile.PhoneNumber,       // 7
		&user.Profile.NIK,               // 8
		&user.Profile.MemberType,        // 9
		&user.Profile.Address,           // 10
		&user.Profile.PhotoKTPURL,       // 11
		&user.Profile.PhotoSelfieURL,    // 12
		&user.Profile.BankName,          // 13
		&user.Profile.BankAccountNumber, // 14
		&user.TotalReferral,             // 15
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Tambahkan parameter `role string` di sini
func (r *UserRepositoryImpl) FindUserByEmail(email string, role string) (*model.UserResponse, error) {
	
	// ==========================================
	// 1. JIKA ROLE ADALAH ADMIN
	// ==========================================
	if role == "admin" {
		queryAdmin := `
			SELECT id, email, full_name, role, is_active
			FROM admins
			WHERE email = ?
			LIMIT 1`
		var admin struct {
			ID       int
			Email    string
			FullName string
			Role     string
			IsActive bool
		}
		
		err := r.DB.QueryRow(queryAdmin, email).Scan(&admin.ID, &admin.Email, &admin.FullName, &admin.Role, &admin.IsActive)
		if err != nil {
			return nil, err // Akan error jika email admin tidak ditemukan
		}
		
		// Cek Shadow Account di tabel users (untuk keperluan FK dan JWT)
		var userID int
		errUser := r.DB.QueryRow("SELECT id FROM users WHERE email = ? AND role = 'admin' LIMIT 1", email).Scan(&userID)
		
		if errUser == sql.ErrNoRows {
			res, errInsert := r.DB.Exec("INSERT INTO users (email, role, status) VALUES (?, 'admin', 'active')", email)
			if errInsert != nil {
				return nil, errInsert
			}
			lastID, errLastID := res.LastInsertId()
			if errLastID != nil {
				return nil, errLastID
			}
			userID = int(lastID)
		} else if errUser != nil {
			return nil, errUser
		}

		status := "active"
		if !admin.IsActive {
			status = "pending"
		}

		return &model.UserResponse{
			ID:     userID, 
			Email:  admin.Email,
			Role:   admin.Role,
			Status: status,
			Profile: model.UserProfile{
				FullName: admin.FullName,
			},
		}, nil
	}

	// ==========================================
	// 2. JIKA ROLE BUKAN ADMIN (USER BIASA)
	// ==========================================
	query := `
		SELECT 
			u.id, u.email, u.role, u.status, u.rejection_reason,
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
		&user.RejectionReason,           // 5 (Pastikan ini sesuai dengan struct kamu)
		&user.Profile.FullName,          // 6
		&user.Profile.PhoneNumber,       // 7
		&user.Profile.NIK,               // 8
		&user.Profile.MemberType,        // 9
		&user.Profile.Address,           // 10
		&user.Profile.PhotoKTPURL,       // 11
		&user.Profile.PhotoSelfieURL,    // 12
		&user.Profile.BankName,          // 13
		&user.Profile.BankAccountNumber, // 14
		&user.TotalReferral,             // 15
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

func (r *UserRepositoryImpl) UpdateRegistrationData(userID int, req model.UpdateRegistrationRequest) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Reset Status di Tabel Users
	queryUser := `
		UPDATE users 
		SET status = 'PENDING', rejection_reason = NULL, verified_at = NULL, verified_by = NULL 
		WHERE id = ?`

	if _, err = tx.Exec(queryUser, userID); err != nil {
		return err
	}

	// 2. Update Data di Tabel User Profiles (Tanpa referral_number)
	queryProfile := `
		UPDATE user_profiles 
		SET full_name = ?, phone_number = ?, nik = ?, member_type = ?, 
		    address = ?, city = ?, bank_name = ?, bank_account_number = ?, 
		    photo_ktp_url = ?, photo_selfie_url = ?
		WHERE user_id = ?`

	_, err = tx.Exec(queryProfile,
		req.FullName, req.PhoneNumber, req.NIK, req.MemberType,
		req.Address, req.City, req.BankName, req.BankAccountNumber,
		req.PhotoKTPURL, req.PhotoSelfieURL, userID,
	)
	if err != nil {
		return err
	}

	// 3. Update Bukti Bayar di Tabel Registration Payments
	queryPayment := `
		UPDATE registration_payments 
		SET payment_proof_url = ? 
		WHERE user_id = ?`

	if _, err = tx.Exec(queryPayment, req.PaymentProofURL, userID); err != nil {
		return err
	}

	// 4. Commit Transaksi
	return tx.Commit()
}
