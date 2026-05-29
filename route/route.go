package route

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/config"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/controller"
	_ "github.com/muhammadfarrasfajri/koperasi-gerai-be/docs"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine, authController *controller.AuthController, userController *controller.UserController, refreshController *controller.RefreshTokenController, walletController *controller.WalletController, jwtManager *middleware.JWTManager, adminController *controller.AdminController) {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/api/auth/v1")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", refreshController.RefreshToken)
		auth.POST("/logout", jwtManager.AuthMiddleware(), authController.Logout)
	}

	// Dev-only helper endpoint to generate JWT token for testing without Firebase.
	// Only active when gin.Mode() == gin.DebugMode.
	r.POST("/api/auth/v1/dev-login", func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.JSON(http.StatusForbidden, gin.H{"error": "dev-login is disabled in release mode"})
			return
		}

		var req struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user in database by email
		var user model.User
		var admin struct {
			ID       int
			Email    string
			Role     string
			IsActive bool
		}
		err := config.DB.QueryRow("SELECT id, email, role, is_active FROM admins WHERE email = ? LIMIT 1", req.Email).Scan(
			&admin.ID, &admin.Email, &admin.Role, &admin.IsActive,
		)
		if err == nil {
			var userID int
			errUser := config.DB.QueryRow("SELECT id FROM users WHERE email = ? AND role = 'admin' LIMIT 1", req.Email).Scan(&userID)
			if errUser == sql.ErrNoRows {
				res, errInsert := config.DB.Exec("INSERT INTO users (email, role, status) VALUES (?, 'admin', 'active')", req.Email)
				if errInsert != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to auto-register admin: " + errInsert.Error()})
					return
				}
				lastID, _ := res.LastInsertId()
				userID = int(lastID)
			} else if errUser != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query admin users: " + errUser.Error()})
				return
			}
			user.ID = userID
			user.Email = admin.Email
			user.Role = admin.Role
			user.Status = "active"
			if !admin.IsActive {
				user.Status = "pending"
			}
		} else {
			err = config.DB.QueryRow("SELECT id, email, role, status FROM users WHERE email = ?", req.Email).Scan(
				&user.ID, &user.Email, &user.Role, &user.Status,
			)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found in local database: " + err.Error()})
				return
			}
		}

		// Generate access token using the global jwtManager
		accessToken, err := jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate refresh token
		refreshToken, err := jwtManager.GenerateRefreshToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Hash and save refresh token to database (Keamanan Utama agar cocok saat refresh)
		refreshTokenHash := middleware.HashToken(refreshToken)
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		sqlQuery := `
			INSERT INTO refresh_tokens (user_id, token, expires_at) 
			VALUES (?, ?, ?) 
			ON DUPLICATE KEY UPDATE token = VALUES(token), expires_at = VALUES(expires_at)`

		_, err = config.DB.Exec(sqlQuery, user.ID, refreshTokenHash, expiresAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	})

	user := r.Group("/api/user/v1")
	{
		user.PATCH("/update-user", jwtManager.AuthMiddleware(), userController.UpdateRegistrationData)
		user.GET("/me", jwtManager.AuthMiddleware(), userController.FindUserByEmail)
		user.GET("/dashboard", jwtManager.AuthMiddleware(), userController.GetUserDashboardDashboard)
	}

	wallet := r.Group("/api/wallet/v1")
	{
		wallet.POST("/withdraw", jwtManager.AuthMiddleware(), walletController.WithdrawAllBalance)
	}

	admin := r.Group("/api/admin/v1")
	admin.Use(jwtManager.AuthMiddleware())
	admin.Use(middleware.AdminOnlyMiddleware())
	{
		admin.GET("/users", adminController.GetUsersList)
		admin.POST("/users/:id/verify", adminController.VerifyRegistration)
		admin.GET("/referrals", adminController.GetReferralTree)
		admin.GET("/dashboard/summary", adminController.GetDashboardSummary)
		admin.GET("/users/:id", adminController.GetUserDetails)
	}
}
