package repository

import (
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AuthRepository interface {
	Register(user model.User, profile model.UserProfile, payment model.RegistrationPayment) error
	Login(rt model.RefreshToken) error
	IsEmailExists(email string) (bool, error)
	IsPhoneNoExists(phoneNumber string) (bool, error)
	IsNIKExists(nik string) (bool, error)
}
