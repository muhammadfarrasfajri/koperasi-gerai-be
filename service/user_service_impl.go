package service

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
)

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

func (s *UserServiceImpl) GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error) {
	return s.UserRepo.GetUserDashboardDashboard(userID)
}
