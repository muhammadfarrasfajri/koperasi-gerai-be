package service

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AuthService interface {
	Register(req model.RegisterMemberRequest) error
	Login(idToken string) (map[string]interface{}, error)
	Logout(rawRefreshToken string) error
}
