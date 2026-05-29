package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash *string   `json:"password_hash"` // Menggunakan pointer (*) agar bisa menerima nilai nil (NULL)
	GoogleID     *string   `json:"google_id"`     // Menggunakan pointer (*) agar bisa menerima nilai nil (NULL)
	AuthProvider string    `json:"auth_provider"` // Contoh isi: "local" atau "google"
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserProfile struct {
	ID                int    `json:"id"`
	UserID            int    `json:"user_id"`
	FullName          string `json:"full_name"`
	PhoneNumber       string `json:"phone_number"`
	NIK               string `json:"nik"`
	MemberType        string `json:"member_type"`
	Address           string `json:"address"`
	PhotoKTPURL       string `json:"photo_ktp_url"`
	PhotoSelfieURL    string `json:"photo_selfie_url"`
	BankName          string `json:"bank_name"`
	BankAccountNumber string `json:"bank_account_number"`
	ReferralNumber    string `json:"referral_number"` // Pointer karena bisa NULL
}

type UserResponse struct {
	ID            int         `json:"id"`
	Email         string      `json:"email"`
	Role          string      `json:"role"`
	IsActive      bool        `json:"is_active"`
	TotalReferral int         `json:"total_referral"`
	Profile       UserProfile `json:"profile"`
}

type RegistrationPayment struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	Amount          float64 `json:"amount"`
	PaymentProofURL string  `json:"payment_proof_url"`
	Status          string  `json:"status"`
}

type RefreshToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterMemberRequest struct {
	TokenId           string  `form:"token_id" binding:"required"`
	FullName          string  `form:"full_name" binding:"required"`
	PhoneNumber       string  `form:"phone_number" binding:"required"`
	NIK               string  `form:"nik" binding:"required"`
	MemberType        string  `form:"member_type" binding:"required"` // perorangan / usaha
	Address           string  `form:"address" binding:"required"`
	BankName          string  `form:"bank_name" binding:"required"`
	BankAccountNumber string  `form:"bank_account_number" binding:"required"`
	ReferralNumber    *string `form:"referral_number"`
	PhotoKTPURL       string  `form:"-"`
	PhotoSelfieURL    string  `form:"-"`
	PaymentProofURL   string  `form:"-"`
}

type LoginMemberRequest struct {
	TokenId string `form:"token_id" binding:"required"`
}
