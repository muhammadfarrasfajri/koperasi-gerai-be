package service

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type AuthService interface {
	Register(user model.User, profile model.UserProfile, payment model.RegistrationPayment) error
}
