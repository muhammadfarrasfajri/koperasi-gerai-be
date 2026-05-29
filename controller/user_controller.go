package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (c *UserController) GetUserDashboardDashboard(ctx *gin.Context) {

	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.APIResponse{
			Error:   true,
			Message: "Unauthorized: Token tidak valid atau user_id tidak ditemukan",
			Type:    "Auth Error",
			Data:    nil,
		})
		return
	}

	result, err := c.UserService.GetUserDashboardDashboard(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "Bad Request",
			Data:    nil,
		})
		return
	}

	// FIX: Wajib mengirimkan response sukses ke client jika proses berhasil!
	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Penarikan berhasil dicatat. Menunggu proses pencairan oleh Admin.",
		Type:    "Success",
		Data:    result,
	})
}

func (c *UserController) UpdateRegistrationData(ctx *gin.Context) {
	var req model.UpdateRegistrationRequest

	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.APIResponse{
			Error:   true,
			Message: "Unauthorized: Token tidak valid atau user_id tidak ditemukan",
			Type:    "Auth Error",
			Data:    nil,
		})
		return
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Invalid data format" + err.Error(),
			Type:    "ValidationError",
		})
		return
	}
	
	// FIX 2: Validasi Foto KTP (Wajib)
	fileKTP, err := ctx.FormFile("photo_ktp_url")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Photo KTP is required",
			Type:    "ValidationError",
		})
		return
	}
	pathKTP := "public/uploads/profile/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileKTP.Filename)
	if errSave := ctx.SaveUploadedFile(fileKTP, pathKTP); errSave != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Error:   true,
			Message: "Failed to save profile photo",
			Type:    "ProfileError",
		})
		return
	}
	req.PhotoKTPURL = pathKTP

	// FIX 3: Validasi Foto Selfie (Wajib)
	fileSelfie, err := ctx.FormFile("photo_selfie_url")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Photo selfie is required",
			Type:    "ValidationError",
		})
		return
	}
	pathSelfie := "public/uploads/profile/" + fmt.Sprintf("%d_%s", time.Now().Unix(), fileSelfie.Filename)
	if errSave := ctx.SaveUploadedFile(fileSelfie, pathSelfie); errSave != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Error:   true,
			Message: "Failed to save selfie photo",
			Type:    "ProfileError",
		})
		return
	}
	req.PhotoSelfieURL = pathSelfie

	// FIX 4: Validasi Bukti Transfer (Wajib)
	filePayment, err := ctx.FormFile("payment_proof_url")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Payment proof photo is required",
			Type:    "ValidationError",
		})
		return
	}
	pathPayment := "public/uploads/payment/" + fmt.Sprintf("%d_%s", time.Now().Unix(), filePayment.Filename)
	if errSave := ctx.SaveUploadedFile(filePayment, pathPayment); errSave != nil {
		ctx.JSON(http.StatusInternalServerError, model.APIResponse{
			Error:   true,
			Message: "Failed to save payment photo",
			Type:    "PaymentError",
		})
		return
	}
	req.PaymentProofURL = pathPayment

	// Panggil Service Layer
	err = c.UserService.UpdateRegistrationData(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "UpdateUserError",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, model.APIResponse{
		Error:   false,
		Message: "update user success",
		Type:    "UpdateUser",
	})
}

func (c *UserController) FindUserByEmail(ctx *gin.Context) {

	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.APIResponse{
			Error:   true,
			Message: "Unauthorized: Token tidak valid atau user_id tidak ditemukan",
			Type:    "Auth Error",
			Data:    nil,
		})
		return
	}

	result, err := c.UserService.FindUserById(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "Bad Request",
			Data:    nil,
		})
		return
	}

	// FIX: Wajib mengirimkan response sukses ke client jika proses berhasil!
	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Berhasil mengembil data",
		Type:    "Success",
		Data:    result,
	})
}
