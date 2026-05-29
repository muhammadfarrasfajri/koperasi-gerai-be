package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterMemberRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Invalid data format" + err.Error(),
			Type:    "ValidationError",
		})
		return
	}

	if req.TokenId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Failed to send id token",
			Type:    "TokenError",
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
	err = c.AuthService.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "RegistrationError",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, model.APIResponse{
		Error:   false,
		Message: "Register Success",
		Type:    "Register user",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginMemberRequest

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Invalid data format",
			Type:    "ValidationError",
		})
		return
	}

	if req.TokenId == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Failed to send id token",
			Type:    "TokenError",
		})
		return
	}

	result, err := c.AuthService.Login(req.TokenId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "LoginError",
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Login successful",
		Type:    "Success",
		Data:    result,
	})

}

func (c *AuthController) Logout(ctx *gin.Context) {

	var req model.RefreshTokenRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Invalid data format",
			Type:    "ValidationError",
		})
		return
	}

	if req.RefreshToken == "" {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: "Failed to send id token",
			Type:    "TokenError",
		})
		return
	}

	err = c.AuthService.Logout(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "Error Logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Logout Success",
		Type:    "Success",
	})

}
