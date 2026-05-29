package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/service"
)

// WebResponse mendefinisikan respons JSON sukses yang terurut
type WebResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// WebErrorResponse mendefinisikan respons JSON error yang terurut
type WebErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type AdminController struct {
	AdminService service.AdminService
}

func NewAdminController(adminService service.AdminService) *AdminController {
	return &AdminController{
		AdminService: adminService,
	}
}

// GetUsersList menangani request GET /api/admin/users?status=xxx&page=1&limit=10
func (ctrl *AdminController) GetUsersList(c *gin.Context) {
	statusFilter := c.Query("status") // 'all', 'pending', 'verified'

	// Ambil parameter paginasi dari query string
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	users, meta, err := ctrl.AdminService.GetUsersList(c.Request.Context(), statusFilter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve users list",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: "Users list retrieved successfully",
		Data:    users,
		Meta:    meta,
	})
}

// VerifyRegistration menangani request POST /api/admin/users/:id/verify
func (ctrl *AdminController) VerifyRegistration(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, WebErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID, must be an integer",
		})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required"` // 'approve' or 'reject'
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, WebErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body, 'action' field is required",
		})
		return
	}

	err = ctrl.AdminService.VerifyRegistration(c.Request.Context(), userID, req.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to verify user registration",
			Details: err.Error(),
		})
		return
	}

	var msg string
	if req.Action == "approve" {
		msg = "User registration approved and user activated successfully"
	} else {
		msg = "User registration rejected successfully"
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: msg,
	})
}

// GetReferralTree menangani request GET /api/admin/referrals
func (ctrl *AdminController) GetReferralTree(c *gin.Context) {
	tree, err := ctrl.AdminService.GetReferralTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to build referral tree",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: "Referral tree retrieved successfully",
		Data:    tree,
	})
}

// GetDashboardSummary menangani request GET /api/admin/v1/dashboard/summary
func (ctrl *AdminController) GetDashboardSummary(c *gin.Context) {
	summary, err := ctrl.AdminService.GetAdminSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve dashboard summary",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: "Dashboard summary retrieved successfully",
		Data:    summary,
	})
}
