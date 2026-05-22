package ad_account

import (
	"net/http"
	"strconv"

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
// @Param        search   query     string  false  "Search by name"
// @Param        page     query     int     false  "Page number"
// @Param        limit    query     int     false  "Items per page"
// @Success      200      {object}  response.Response{data=[]dto.AdAccountResponse,meta=response.Meta}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts [get]
func (h *Handler) GetAdAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	search := c.Query("search")

	filter := AdAccountFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	resp, meta, err := h.service.GetAdAccounts(filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithMeta(c, "Successfully retrieved Meta ad accounts", resp, *meta)
}
