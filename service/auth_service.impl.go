package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
)

type AuthServiceImpl struct {
	AuthRepo     repository.AuthRepository
	UserRepo     repository.UserRepository
	RefreshRepo  repository.RefreshTokenRepository
	JWTManager   *middleware.JWTManager
	FirebaseAuth *auth.Client
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, refreshRepo repository.RefreshTokenRepository, jwtManger *middleware.JWTManager, userAuth *auth.Client) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepo:     authRepo,
		UserRepo:     userRepo,
		RefreshRepo:  refreshRepo,
		JWTManager:   jwtManger,
		FirebaseAuth: userAuth,
	}
}

var numericRegex = regexp.MustCompile(`^[0-9]+$`)

func (s *AuthServiceImpl) Register(req model.RegisterMemberRequest) error {
	var paymentAmount float64 = 1000000.00

	// 1. Verifikasi Firebase Token
	token, err := s.FirebaseAuth.VerifyIDToken(context.Background(), req.TokenId)
	if err != nil {
		return fmt.Errorf("invalid google token: %v", err)
	}

	email, _ := token.Claims["email"].(string)

	// 2. Pembersihan & Validasi Format Data
	cleanPhone := strings.TrimSpace(req.PhoneNumber)
	panjangPhone := len(cleanPhone)
	if panjangPhone < 10 || panjangPhone > 15 {
		return errors.New("invalid phone number length")
	}
	if !numericRegex.MatchString(cleanPhone) {
		return errors.New("phone number must contain only digits")
	}

	cleanNIK := strings.TrimSpace(req.NIK)
	if len(cleanNIK) != 16 {
		return errors.New("invalid NIK length, must be exactly 16 digits")
	}
	if !numericRegex.MatchString(cleanNIK) {
		return errors.New("NIK must contain only digits")
	}
	if req.ReferralNumber == "" {
		return errors.New("referral must be filled in")
	}

	// 3. Validasi Data Unik ke Database
	isExistsEmail, err := s.AuthRepo.IsEmailExists(email)
	if err != nil {
		return err
	}
	if isExistsEmail {
		return fmt.Errorf("email already exists")
	}

	isExistsPhone, err := s.AuthRepo.IsPhoneNoExists(cleanPhone)
	if err != nil {
		return err
	}
	if isExistsPhone {
		return fmt.Errorf("phone number already exists")
	}

	isExistsNIK, err := s.AuthRepo.IsNIKExists(cleanNIK)
	if err != nil {
		return err
	}
	if isExistsNIK {
		return fmt.Errorf("NIK already exists")
	}

	// 4. Mapping DTO ke Model User
	user := model.User{
		Email:    email,
		GoogleID: &token.UID,
		Role:     "member",
		IsActive: false,
	}

	// 5. Mapping DTO ke Model Profile
	profile := model.UserProfile{
		FullName:          req.FullName,
		PhoneNumber:       cleanPhone,
		NIK:               cleanNIK,
		MemberType:        req.MemberType,
		Address:           req.Address,
		PhotoKTPURL:       req.PhotoKTPURL,
		PhotoSelfieURL:    req.PhotoSelfieURL,
		BankName:          req.BankName,
		BankAccountNumber: req.BankAccountNumber,
		ReferralNumber:    req.ReferralNumber,
	}

	// 6. Mapping Pembayaran
	// FIX: Menggunakan model. bukan dto.
	payment := model.RegistrationPayment{
		Amount:          paymentAmount,
		PaymentProofURL: req.PaymentProofURL,
	}

	// 8. Kirim ke Repository (Database Transaction)
	// FIX: Menambahkan parameter emptyRefresh di paling akhir
	err = s.AuthRepo.Register(user, profile, payment)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) Login(idToken string) (map[string]interface{}, error) {
	ctx := context.Background()

	// verify idToken with Firebase Auth
	token, err := s.FirebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("invalid ID token: " + err.Error())
	}

	// Extract email from token claims
	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("email not found in token claims")
	}

	// find user by email
	user, err := s.UserRepo.FindUserByEmail(email)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("user not found, please register first")
	}

	// generate access token
	accessToken, err := s.JWTManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, err := s.JWTManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Hash refresh token sebelum disimpan di database
	refreshTokenHash := middleware.HashToken(refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	tokenData := model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenHash,
		ExpiresAt: expiresAt,
	}

	// Upsert token login (Update jika sudah ada, Insert jika belum)
	err = s.RefreshRepo.UpsertTokenLogin(tokenData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Profile.FullName,
			"email": user.Email,
			"role":  user.Role,
		},
	}, nil
}

func (s *AuthServiceImpl) Logout(rawRefreshToken string) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// 1. LOG: Awal Proses
	// Kita log hash-nya saja, jangan raw token-nya demi keamanan
	fmt.Printf("[AUTH-LOGOUT] [%s] Memulai proses logout...\n", now)

	tokenHash := middleware.HashToken(rawRefreshToken)

	// 2. Proses hapus di Repository
	err := s.RefreshRepo.DeleteRefreshToken(tokenHash)
	if err != nil {
		// LOG: Jika gagal hapus di DB
		fmt.Printf("[ERROR] [%s] Gagal menghapus refresh token saat logout: %v\n", now, err)
		return errors.New("logout failed")
	}

	// 3. LOG: Berhasil
	fmt.Printf("[SUCCESS] [%s] User berhasil logout. Token hash: %s...\n", now, tokenHash[:10])
	return nil
}
