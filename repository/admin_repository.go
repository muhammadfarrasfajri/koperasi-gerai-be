package repository

import (
	"context"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AdminRepository interface {
	GetUsers(ctx context.Context, statusFilter string, page int, limit int) ([]model.UserDetail, int64, error)
	VerifyUser(ctx context.Context, userID int64, paymentStatus string, userStatus string, rejectionReason string, verifiedBy int64) error
	GetAllProfiles(ctx context.Context) ([]model.UserProfile, error)
	GetAdminSummary(ctx context.Context) (model.AdminSummary, error)
	GetUserByID(ctx context.Context, userID int64) (model.UserDetail, error)
}
