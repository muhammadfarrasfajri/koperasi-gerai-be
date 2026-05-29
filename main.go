package main

import (
	"log"
	"os" // Diperlukan untuk membaca os.Getenv

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/bootstrap"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/middleware"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/route"
)

// @title          Koperasi Gerai API
// @version        1.0
// @description    Backend API Service untuk Aplikasi Koperasi Gerai.
// @termsOfService http://swagger.io/terms/

// @contact.name  Muhammad Farras Fajri
// @contact.email farrasfajri@example.com

// @host      localhost:8080
// @BasePath  /
func main() {
	err := godotenv.Load()
	if err != nil {
		// Log fatal akan menghentikan aplikasi jika file .env tidak ditemukan
		log.Fatal("Error loading .env file")
	}

	// 1. Inisialisasi Database
	bootstrap.InitDatabase()

	// 2. Inisialisasi Firebase (untuk Google Auth)
	firebaseAuth := bootstrap.InitFirebase()

	// 3. Inisialisasi Dependency Injection (Container)
	container := bootstrap.InitContainer(firebaseAuth)

	// 4. Inisialisasi Gin Engine
	r := gin.Default()

	// 6. Global Middleware
	middleware.AttachCORS(r)

	// 5. Konfigurasi Static Files & Upload Memory
	// Mengizinkan akses file di folder public (seperti foto profil)
	r.Static("/public", "./public")
	// Limit memory untuk upload file (50 MiB)
	r.MaxMultipartMemory = 50 << 20

	// 7. Setup Routes
	// Melewatkan semua controller dari container ke router
	route.SetupRouter(
		r,
		container.AuthController,
		container.UserController,
		container.RefreshController,
		container.WalletController,
		container.JWTManager,
		container.AdminController,
	)

	// 8. Jalankan Server (FIX: Membaca port dari .env, fallback ke 8080 jika kosong)
	port := os.Getenv("PORT")

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
