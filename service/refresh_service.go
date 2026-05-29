package service

type RefreshTokenService interface {
	RefreshToken(refreshToken string) (map[string]interface{}, error)
}
