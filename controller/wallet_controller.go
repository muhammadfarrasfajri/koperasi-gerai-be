package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

type WalletController struct {
	WalletService service.WalletService
}

func NewWalletController(walletService service.WalletService) *WalletController {
	return &WalletController{
		WalletService: walletService,
	}
}

func (c *WalletController) WithdrawAllBalance(ctx *gin.Context) {

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

	err := c.WalletService.WithdrawAllBalance(userID)
	if err != nil {
		// FIX: Ubah status menjadi 400 Bad Request dan tampilkan pesan error asli dari service/repo
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(), // Akan menampilkan "saldo kosong, tidak ada yang bisa ditarik" dll
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
		Data:    nil,
	})
}
