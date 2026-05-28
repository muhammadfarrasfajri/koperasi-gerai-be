package controller

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/service"

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}
