package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type AuthRepository interface {
	Register(user model.User, profile model.UserProfile, payment model.RegistrationPayment) error
}