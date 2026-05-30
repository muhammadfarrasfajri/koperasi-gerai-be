package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
)

type AdminServiceImpl struct {
	AdminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) *AdminServiceImpl {
	return &AdminServiceImpl{
		AdminRepo: adminRepo,
	}
}

func (s *AdminServiceImpl) GetUsersList(ctx context.Context, statusFilter string, page int, limit int) ([]model.UserDetail, model.PaginationMetadata, error) {
	// 1. Tentukan nilai default jika parameter tidak valid
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// 2. Tarik data terpaginasi dari database
	users, totalRows, err := s.AdminRepo.GetUsers(ctx, statusFilter, page, limit)
	if err != nil {
		return nil, model.PaginationMetadata{}, err
	}

	// 3. Hitung total halaman
	totalPages := int(totalRows / int64(limit))
	if totalRows%int64(limit) > 0 {
		totalPages++
	}

	// 4. Susun metadata paginasi
	meta := model.PaginationMetadata{
		CurrentPage: page,
		Limit:       limit,
		TotalRows:   totalRows,
		TotalPages:  totalPages,
	}

	return users, meta, nil
}

func (s *AdminServiceImpl) VerifyRegistration(ctx context.Context, userID int64, action string, rejectionReason string, verifiedBy int64) error {
	var paymentStatus string
	var userStatus string

	switch action {
	case "approve", "active":
		paymentStatus = "approve"
		userStatus = "active"
	case "reject":
		paymentStatus = "reject"
		userStatus = "reject"
	default:
		return errors.New("invalid action: action must be 'approve', 'active' or 'reject'")
	}

	return s.AdminRepo.VerifyUser(ctx, userID, paymentStatus, userStatus, rejectionReason, verifiedBy)
}

func (s *AdminServiceImpl) GetReferralTree(ctx context.Context) ([]*model.ReferralNode, error) {
	profiles, err := s.AdminRepo.GetAllProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles for referral tree: %w", err)
	}

	return BuildReferralTree(profiles), nil
}

func (s *AdminServiceImpl) GetAdminSummary(ctx context.Context) (model.AdminSummary, error) {
	return s.AdminRepo.GetAdminSummary(ctx)
}

func (s *AdminServiceImpl) GetUserDetails(ctx context.Context, userID int64) (model.UserDetail, error) {
	return s.AdminRepo.GetUserByID(ctx, userID)
}

// BuildReferralTree mengonstruksi pohon referral berkedalaman N secara efisien O(N) menggunakan pointer map
func BuildReferralTree(profiles []model.UserProfile) []*model.ReferralNode {
	nodesMap := make(map[string]*model.ReferralNode)

	// Step 1: Daftarkan semua profile ke dalam map lookup
	for _, p := range profiles {
		whatsappNum := p.PhoneNumber
		var refNumPtr *string
		if p.ReferralNumber != "" {
			refVal := p.ReferralNumber
			refNumPtr = &refVal
		}
		nodesMap[whatsappNum] = &model.ReferralNode{
			UserID:         int64(p.UserID),
			FullName:       p.FullName,
			WhatsAppNumber: whatsappNum,
			ReferralNumber: refNumPtr,
			Children:       []*model.ReferralNode{},
		}
	}

	var roots []*model.ReferralNode

	// Step 2: Hubungkan simpul anak ke simpul orang tua
	for _, node := range nodesMap {
		if node.ReferralNumber == nil || *node.ReferralNumber == "" {
			// Simpul utama (tidak memiliki referral)
			roots = append(roots, node)
		} else {
			parent, exists := nodesMap[*node.ReferralNumber]
			if exists {
				// Cegah self-reference sederhana
				if parent.WhatsAppNumber == node.WhatsAppNumber {
					roots = append(roots, node)
				} else {
					parent.Children = append(parent.Children, node)
				}
			} else {
				// Jika nomor referral orang tua tidak terdaftar di database, anggap sebagai root
				roots = append(roots, node)
			}
		}
	}

	return roots
}

func (s *AdminServiceImpl) VerifyWithdrawal(adminID int, req model.VerifyWithdrawalRequest) error {
	// A. Standarisasi input menjadi huruf kecil (lowercase)
	// Ini mencegah bug jika frontend mengirim "REJECT", "Reject", atau "approve"
	req.Status = strings.ToLower(strings.TrimSpace(req.Status))

	// B. Validasi aksi yang diperbolehkan
	if req.Status != "approve" && req.Status != "reject" {
		return errors.New("status tidak valid: hanya menerima 'approve' atau 'reject'")
	}

	// C. Aturan Bisnis: Wajib isi alasan jika ditolak
	if req.Status == "reject" {
		if req.RejectReason == nil || strings.TrimSpace(*req.RejectReason) == "" {
			return errors.New("alasan penolakan (reject_reason) wajib diisi jika menolak pencairan")
		}
	} else if req.Status == "approve" {
		// Bersihkan data alasan jika disetujui, agar tabel database tetap bersih (NULL)
		req.RejectReason = nil
	}

	// D. Oper data yang sudah tervalidasi dan bersih ke layer Repository
	return s.AdminRepo.VerifyWithdrawal(adminID, req)
}
