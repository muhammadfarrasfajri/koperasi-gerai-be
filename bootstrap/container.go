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
	UserController    *controller.UserController
	RefreshController *controller.RefreshTokenController
	WalletController  *controller.WalletController
	JWTManager        *middleware.JWTManager
	AdminController   *controller.AdminController // Integrasi Admin Controller
}

func InitContainer(userAuth *auth.Client) *Container {

	// Initialize repositories
	authRepo := repository.NewAuthRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	refreshRepo := repository.NewRefreshTokenRepository(config.DB)
	walletRepo := repository.NewWalletRepositoryImpl(config.DB)
	adminRepo := repository.NewAdminRepository() // Inisialisasi Admin GORM Repo

	// Initialize middleware
	jwtManager := middleware.NewJWTManager(os.Getenv("ACCESS_TOKEN_SECRET"), os.Getenv("REFRESH_TOKEN_SECRET"))

	// Initialize services
	authService := service.NewAuthService(authRepo, userRepo, refreshRepo, jwtManager, userAuth)
	userService := service.NewUserService(userRepo)
	refreshService := service.NewRefreshTokenService(refreshRepo, userRepo, jwtManager)
	walletService := service.NewWalletService(walletRepo)
	adminService := service.NewAdminService(adminRepo) // Inisialisasi Admin Service

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	refreshController := controller.NewRefreshTokenController(refreshService)
	walletController := controller.NewWalletController(walletService)
	adminController := controller.NewAdminController(adminService) // Inisialisasi Admin Controller

	return &Container{
		AuthController:    authController,
		UserController:    userController,
		RefreshController: refreshController,
		WalletController:  walletController,
		JWTManager:        jwtManager,
		AdminController:   adminController, // Kembalikan Admin Controller
	}
}
