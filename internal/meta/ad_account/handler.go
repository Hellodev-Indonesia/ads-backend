package ad_account

import (
	"net/http"

	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.AdAccountResponse{}

// GetAdAccounts godoc
// @Summary      Get Meta Ad Accounts
// @Description  Retrieve all ad accounts associated with the system user token
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Success      200      {object}  response.SuccessResponse{data=[]dto.AdAccountResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts [get]
func (h *Handler) GetAdAccounts(c *gin.Context) {
	resp, err := h.service.GetAdAccounts()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved Meta ad accounts", resp)
}
