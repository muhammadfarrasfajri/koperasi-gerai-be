package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type UserRepository interface {
	FindUserById(id int) (*model.UserResponse, error)
	FindUserByEmail(email string) (*model.UserResponse, error)
	GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error)
}
