package service

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"

type AuthServiceImpl struct {
	AuthRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepo: authRepo,
	}
}
