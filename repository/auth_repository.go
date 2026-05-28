package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type AuthRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

func (r *AuthRepositoryImpl) Register(user model.User, profile model.UserProfile, payment model.RegistrationPayment) error {
	return errors.New("")
}
