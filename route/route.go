package route

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/controller"
	_ "github.com/muhammadfarrasfajri/koperasi-gerai-be/docs"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine, authController *controller.AuthController, userController *controller.UserController, refreshController *controller.RefreshTokenController, walletController *controller.WalletController, jwtManager *middleware.JWTManager) {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/api/auth/v1")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", refreshController.RefreshToken)
		auth.POST("logout", jwtManager.AuthMiddleware(), authController.Logout)
	}

	user := r.Group("/api/user/v1")
	{
		user.GET("/dashboard", jwtManager.AuthMiddleware(), userController.GetUserDashboardDashboard)
	}

	wallet := r.Group("/api/wallet/v1")
	{
		wallet.POST("/withdraw", jwtManager.AuthMiddleware(), walletController.WithdrawAllBalance)
	}
}
