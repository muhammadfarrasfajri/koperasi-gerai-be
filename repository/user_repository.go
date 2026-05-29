package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type UserRepository interface {
	FindUserById(id int, role string) (*model.UserResponse, error)
	FindUserByEmail(email string, role string) (*model.UserResponse, error)
	GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error)
	UpdateRegistrationData(userID int, req model.UpdateRegistrationRequest) error
}
