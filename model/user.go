package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserProfile struct {
	ID                int64   `json:"id"`
	UserID            int64   `json:"user_id"`
	FullName          string  `json:"full_name"`
	WhatsAppNumber    string  `json:"whatsapp_number"`
	NIK               string  `json:"nik"`
	MemberType        string  `json:"member_type"`
	Address           string  `json:"address"`
	PhotoKTPURL       string  `json:"photo_ktp_url"`
	PhotoSelfieURL    string  `json:"photo_selfie_url"`
	BankName          string  `json:"bank_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	ReferralNumber    *string `json:"referral_number"` // Pointer karena bisa NULL
}

type RegistrationPayment struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	Amount          float64 `json:"amount"`
	PaymentProofURL string  `json:"payment_proof_url"`
	Status          string  `json:"status"`
}
