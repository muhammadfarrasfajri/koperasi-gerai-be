package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

func NewJWTManager(accessSecret, refreshSecret string) *JWTManager {
	return &JWTManager{
		AccessSecret:  []byte(accessSecret),
		RefreshSecret: []byte(refreshSecret),
	}
}

// Generate JWT Token (Menggunakan int64 lebih aman untuk database ID)
func (j *JWTManager) GenerateAccessToken(userID int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessSecret)
}

func (j *JWTManager) GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshSecret)
}

// Middleware untuk validasi token
func (j *JWTManager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Format("15:04:05")
		path := c.Request.URL.Path

		fmt.Printf("[AUTH-GATE] [%s] Memeriksa akses untuk path: %s\n", now, path)

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Printf("[AUTH-DENIED] [%s] Header kosong atau tidak pakai 'Bearer '\n", now)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return j.AccessSecret, nil // FIX: Langsung return j.AccessSecret tanpa cast []byte lagi karena tipenya sudah []byte
		})

		if err != nil || !token.Valid {
			fmt.Printf("[AUTH-DENIED] [%s] Token tidak valid/expired: %v\n", now, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("[AUTH-DENIED] [%s] Gagal membaca MapClaims\n", now)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// FIX 1: Ambil sebagai float64 terlebih dahulu, lalu konversi ke int64/int
		userIDFloat, okUserID := claims["user_id"].(float64)
		if !okUserID {
			fmt.Printf("[AUTH-DENIED] [%s] Token valid tapi 'user_id' tidak ditemukan atau bukan angka!\n", now)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id is missing in token"})
			c.Abort()
			return
		}
		userID := int(userIDFloat)

		// Set data ke context dan izinkan masuk
		c.Set("user_id", userID)

		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		// FIX 2: Mengubah %s menjadi %d karena userID bertipe int64 (angka)
		fmt.Printf("[AUTH-OK] [%s] Akses diterima. UserID: %d, Path: %s\n", now, userID, path)
		c.Next()
	}
}

// AdminOnlyMiddleware memastikan pengguna memiliki role 'admin' di context Gin
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Access forbidden: only admin role is allowed",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
