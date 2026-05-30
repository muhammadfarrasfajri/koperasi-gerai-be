package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muhammadfarrasfajri/koperasi-gerai-be/model"
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
		Action          string `json:"action" binding:"required"` // 'approve', 'active', or 'reject'
		RejectionReason string `json:"rejection_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, WebErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body, 'action' field is required",
		})
		return
	}

	if req.Action == "reject" && req.RejectionReason == "" {
		c.JSON(http.StatusBadRequest, WebErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "rejection_reason is required when action is reject",
		})
		return
	}

	// Ambil ID admin yang memverifikasi dari context (diekstrak oleh AuthMiddleware)
	adminIDVal, exists := c.Get("user_id")
	var adminID int64
	if exists {
		if idInt, ok := adminIDVal.(int); ok {
			adminID = int64(idInt)
		} else if idFloat, ok := adminIDVal.(float64); ok {
			adminID = int64(idFloat)
		} else if idInt64, ok := adminIDVal.(int64); ok {
			adminID = idInt64
		}
	}

	err = ctrl.AdminService.VerifyRegistration(c.Request.Context(), userID, req.Action, req.RejectionReason, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, WebErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to verify user registration",
			Details: err.Error(),
		})
		return
	}

	var msg string
	if req.Action == "approve" || req.Action == "active" {
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

// GetUserDetails menangani request GET /api/admin/v1/users/:id
func (ctrl *AdminController) GetUserDetails(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, WebErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID, must be an integer",
		})
		return
	}

	user, err := ctrl.AdminService.GetUserDetails(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, WebErrorResponse{
			Status:  http.StatusNotFound,
			Message: "User not found or is not a member",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: "User details retrieved successfully",
		Data:    user,
	})
}

func (c *AdminController) VerifyWithdrawal(ctx *gin.Context) {
	// A. Ambil ID Admin dari JWT Middleware
	adminID := ctx.GetInt("user_id")
	if adminID == 0 {
		ctx.JSON(http.StatusUnauthorized, WebResponse{
			Status:  http.StatusUnauthorized,
			Message: "Akses ditolak: Token admin tidak valid atau tidak ditemukan",
			Data:    nil,
		})
		return
	}

	// B. Tangkap data JSON dari Frontend
	var req model.VerifyWithdrawalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, WebResponse{
			Status:  http.StatusBadRequest,
			Message: "Format data tidak valid: pastikan semua kolom wajib diisi",
			Data:    nil,
		})
		return
	}

	// C. Oper ke layer Service
	err := c.AdminService.VerifyWithdrawal(adminID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, WebResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// D. Siapkan pesan sukses dinamis
	pesanSukses := "Pencairan referral berhasil disetujui, dana bisa segera ditransfer"
	if strings.ToLower(req.Status) == "reject" {
		pesanSukses = "Pencairan berhasil ditolak. Saldo referral telah dikembalikan ke pengguna"
	}

	// E. Kirim respons sukses dengan WebResponse
	ctx.JSON(http.StatusOK, WebResponse{
		Status:  http.StatusOK,
		Message: pesanSukses,
		Data:    nil, // Bisa diisi nil karena tidak ada data spesifik yang perlu dikembalikan setelah update
	})
}
