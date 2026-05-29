package bootstrap

import (
	"os"

	"firebase.google.com/go/auth"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/config"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/controller"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

type Container struct {
	AuthController    *controller.AuthController
	RefreshController *controller.RefreshTokenController
	WalletController  *controller.WalletController
	JWTManager        *middleware.JWTManager
}

func InitContainer(userAuth *auth.Client) *Container {

	// Initialize repositories
	authRepo := repository.NewAuthRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	refreshRepo := repository.NewRefreshTokenRepository(config.DB)
	walletRepo := repository.NewWalletRepositoryImpl(config.DB)

	// Initialize middleware
	jwtManager := middleware.NewJWTManager(os.Getenv("ACCESS_TOKEN_SECRET"), os.Getenv("REFRESH_TOKEN_SECRET"))

	// Initialize services
	authService := service.NewAuthService(authRepo, userRepo, refreshRepo, jwtManager, userAuth)
	refreshService := service.NewRefreshTokenService(refreshRepo, userRepo, jwtManager)
	walletService := service.NewWalletService(walletRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	refreshController := controller.NewRefreshTokenController(refreshService)
	walletController := controller.NewWalletController(walletService)

	return &Container{
		AuthController:    authController,
		RefreshController: refreshController,
		WalletController:  walletController,
		JWTManager:        jwtManager,
	}
}
