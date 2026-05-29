package service

import (
	"errors"
	"strings"

	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/repository"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/utils"
)

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

func (s *UserServiceImpl) GetUserDashboardDashboard(userID int) (*model.UserDashboardResponse, error) {
	return s.UserRepo.GetUserDashboardDashboard(userID)
}

func (s *UserServiceImpl) FindUserById(id int, role string) (*model.UserResponse, error) {
	return s.UserRepo.FindUserById(id, role)
}

func (s *UserServiceImpl) UpdateRegistrationData(userID int, req model.UpdateRegistrationRequest) error {
	// 1. Ekstrak Kota dari NIK menggunakan fungsi utils yang sudah kita buat
	cityName, err := utils.GetCityFromNIK(req.NIK)
	if err != nil {
		return errors.New("NIK tidak valid: " + err.Error())
	}

	// Masukkan hasil ekstraksi kota ke dalam request untuk disimpan ke database
	req.City = cityName

	// 2. Validasi Tipe Anggota (Business Logic)
	// Memastikan data yang masuk sesuai dengan standar database
	memberType := strings.ToLower(req.MemberType)
	if memberType != "perorangan" && memberType != "usaha" {
		return errors.New("tipe anggota tidak valid, harus 'perorangan' atau 'usaha'")
	}
	req.MemberType = memberType // Format ulang agar selalu huruf kecil jika diperlukan

	// 3. Panggil layer Repository
	// Catatan: Karena di Controller `userID` bertipe `int` (dari ctx.GetInt),
	// kita lakukan casting menjadi `int64` saat memanggil Repository
	err = s.UserRepo.UpdateRegistrationData(userID, req)
	if err != nil {
		return errors.New("terjadi kesalahan saat memperbarui data di database: " + err.Error())
	}

	return nil
}
