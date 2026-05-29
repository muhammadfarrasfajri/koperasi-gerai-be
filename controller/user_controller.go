package controller

import (
	"net/http"

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
