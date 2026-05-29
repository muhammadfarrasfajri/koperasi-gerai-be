package repository

import (
	"database/sql"
	"errors"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
)

type RefreshTokenRepositoryImpl struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepositoryImpl {
	return &RefreshTokenRepositoryImpl{
		DB: db,
	}
}

func (r *RefreshTokenRepositoryImpl) FindRefreshTokenUser(userID int) (*model.RefreshToken, error) {
	sqlQuery := `
		SELECT user_id, token, expires_at
		FROM refresh_tokens
		WHERE user_id = ?
		LIMIT 1
	`

	row := r.DB.QueryRow(sqlQuery, userID)

	token := model.RefreshToken{}
	err := row.Scan(&token.ID, &token.Token, &token.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepositoryImpl) UpsertTokenLogin(rt model.RefreshToken) error {
	// Pakai ON DUPLICATE KEY UPDATE agar token lama tertimpa token baru
	sqlQuery := `
        INSERT INTO refresh_tokens (user_id, token, expires_at) 
        VALUES (?, ?, ?) 
        ON DUPLICATE KEY UPDATE token = VALUES(token), expires_at = VALUES(expires_at)`

	_, err := r.DB.Exec(sqlQuery, rt.UserID, rt.Token, rt.ExpiresAt)
	return err
}

func (r *RefreshTokenRepositoryImpl) DeleteRefreshToken(token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = ?"

	result, err := r.DB.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}
