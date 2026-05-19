package auth

import (
	"net/http"

	"github.com/alex/ads_backend/internal/core/auth/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Login Request"
// @Success      200      {object}  response.SuccessResponse{data=dto.LoginResponse}
// @Failure      401      {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, "Login successful", resp)
}

// Logout godoc
// @Summary      Logout
// @Description  Logout user (Clear session/token)
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.SuccessResponse
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// Logic to invalidate token can be added here (e.g., Redis blacklist)
	response.Success(c, "Logout successful", nil)
}
