package service

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type UserService interface {
	GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error)
}
