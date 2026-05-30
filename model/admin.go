package model

import "time"

// UserDetail mengombinasikan data User, UserProfile, dan RegistrationPayment untuk Admin View
type UserDetail struct {
	UserID            int64     `json:"user_id"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	FullName          string    `json:"full_name"`
	WhatsAppNumber    string    `json:"whatsapp_number" gorm:"column:whatsapp_number"`
	NIK               string    `json:"nik"`
	MemberType        string    `json:"member_type"`
	Address           string    `json:"address"`
	City              string    `json:"city"`
	PhotoKTPURL       string    `json:"photo_ktp_url"`
	PhotoSelfieURL    string    `json:"photo_selfie_url"`
	BankName          string    `json:"bank_name"`
	BankAccountNumber string    `json:"bank_account_number"`
	ReferralNumber    *string   `json:"referral_number"`
	PaymentAmount     float64   `json:"payment_amount"`
	PaymentProofURL   string    `json:"payment_proof_url"`
	RejectionReason   string    `json:"rejection_reason"`
	ReferredByID      *int64    `json:"referred_by_id"`
}

// ReferralNode merepresentasikan satu simpul dalam pohon referral
type ReferralNode struct {
	UserID         int64           `json:"user_id"`
	FullName       string          `json:"full_name"`
	WhatsAppNumber string          `json:"whatsapp_number"`
	ReferralNumber *string         `json:"referral_number"`
	Children       []*ReferralNode `json:"children"`
}

// PaginationMetadata menyimpan informasi metadata halaman untuk daftar data
type PaginationMetadata struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalRows   int64 `json:"total_rows"`
	TotalPages  int   `json:"total_pages"`
}

// AdminSummary menyimpan ringkasan data statistik untuk Dashboard Admin
type AdminSummary struct {
	TotalUsers               int64 `json:"total_users"`
	TotalPendingVerification int64 `json:"total_pending_verifications"`
}
