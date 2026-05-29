package repository

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/model"

type RefreshTokenRepository interface {
	FindRefreshTokenUser(userID int) (*model.RefreshToken, error)
	UpsertTokenLogin(rt model.RefreshToken) error
	DeleteRefreshToken(token string) error
}
