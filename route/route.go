package route

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/controller"
	_ "github.com/muhammadfarrasfajri/koperasi-gerai-be/docs"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
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
		admin.POST("/referral/verify-withdrawal", adminController.VerifyWithdrawal)
	}
}
