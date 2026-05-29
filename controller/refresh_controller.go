package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

type RefreshTokenController struct {
	RefreshService service.RefreshTokenService
}

func NewRefreshTokenController(refreshService service.RefreshTokenService) *RefreshTokenController {
	return &RefreshTokenController{
		RefreshService: refreshService,
	}
}

func (c *RefreshTokenController) RefreshToken(ctx *gin.Context) {
	var req model.RefreshTokenRequest
	err := ctx.ShouldBind(&req)
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

	result, err := c.RefreshService.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.APIResponse{
			Error:   true,
			Message: err.Error(),
			Type:    "Error refresh token",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.APIResponse{
		Error:   false,
		Message: "Register Success",
		Type:    "Success",
		Data:    result,
	})

}
