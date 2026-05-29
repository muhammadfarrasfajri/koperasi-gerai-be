package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
)

type RefreshTokenServiceImpl struct {
	RefreshRepo repository.RefreshTokenRepository
	UserRepo    repository.UserRepository
	JWTManager  *middleware.JWTManager
}

func NewRefreshTokenService(refreshRepo repository.RefreshTokenRepository, userRepo repository.UserRepository, jwtManager *middleware.JWTManager) *RefreshTokenServiceImpl {
	return &RefreshTokenServiceImpl{
		RefreshRepo: refreshRepo,
		UserRepo:    userRepo,
		JWTManager:  jwtManager,
	}
}

func (s *RefreshTokenServiceImpl) RefreshToken(refreshToken string) (map[string]interface{}, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	// 1. LOG: Memulai proses Refresh
	fmt.Printf("[AUTH-REFRESH] [%s] Mencoba rotasi token...\n", now)

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error JWT")
		}
		return s.JWTManager.RefreshSecret, nil
	})

	if err != nil {
		fmt.Printf("[ERROR] [%s] Gagal parsing refresh token: %v\n", now, err)
		return nil, errors.New("Error parsing refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Printf("[WARN] [%s] Token claims tidak valid atau token sudah expired\n", now)
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		fmt.Printf("[WARN] [%s] Token valid tapi user_id tidak ditemukan di claims\n", now)
		return nil, errors.New("invalid user id in token")
	}
	userID := int(userIDFloat)

	// 2. LOG: Cari di Database
	tokenCheck, err := s.RefreshRepo.FindRefreshTokenUser(userID)
	if err != nil {
		return nil, err
	}

	// 3. LOG: Validasi Hash (Keamanan Utama)
	incomingTokenHash := middleware.HashToken(refreshToken)
	if incomingTokenHash != tokenCheck.Token {
		return nil, errors.New("refresh token not match")
	}

	user, err := s.UserRepo.FindUserById(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// 4. LOG: Generate Token Baru
	newAccessToken, err := s.JWTManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		fmt.Printf("[ERROR] [%s] Gagal generate Access Token baru: %v\n", now, err)
		return nil, errors.New("Error generating access token")
	}

	newRefreshToken, err := s.JWTManager.GenerateRefreshToken(user.ID)
	if err != nil {
		fmt.Printf("[ERROR] [%s] Gagal generate Refresh Token baru: %v\n", now, err)
		return nil, errors.New("Error generating refresh token")
	}

	newRefreshTokenHash := middleware.HashToken(newRefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	tokenModel := model.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshTokenHash,
		ExpiresAt: expiresAt,
	}

	// 5. LOG: Update Session di DB
	err = s.RefreshRepo.UpsertTokenLogin(tokenModel)
	if err != nil {
		fmt.Printf("[ERROR] [%s] Gagal menyimpan session baru ke database: %v\n", now, err)
		return nil, errors.New("Failed to save new session")
	}

	return map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, nil
}
