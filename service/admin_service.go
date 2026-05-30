package service

import (
	"context"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AdminService interface {
	GetUsersList(ctx context.Context, statusFilter string, page int, limit int) ([]model.UserDetail, model.PaginationMetadata, error)
	VerifyRegistration(ctx context.Context, userID int64, action string, rejectionReason string, verifiedBy int64) error
	GetReferralTree(ctx context.Context) ([]*model.ReferralNode, error)
	GetAdminSummary(ctx context.Context) (model.AdminSummary, error)
	GetUserDetails(ctx context.Context, userID int64) (model.UserDetail, error)
	VerifyWithdrawal(adminID int, req model.VerifyWithdrawalRequest) error
}
