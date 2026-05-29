package service

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type UserService interface {
	GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error)
	UpdateRegistrationData(userID int, req model.UpdateRegistrationRequest) error
	FindUserById(id int, role string) (*model.UserResponse, error)
}
