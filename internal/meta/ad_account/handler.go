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
// @Param        fields   query     string  false  "Custom fields comma-separated preset"
// @Param        limit    query     string  false  "Pagination limit"
// @Param        after    query     string  false  "Cursor after"
// @Param        before   query     string  false  "Cursor before"
// @Success      200      {object}  response.Response{data=[]dto.AdAccountResponse,paging=response.MetaPaging}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts [get]
func (h *Handler) GetAdAccounts(c *gin.Context) {
	fields := c.Query("fields")
	if fields == "" {
		fields = DefaultFields
	}

	limit := c.Query("limit")
	after := c.Query("after")
	before := c.Query("before")

	resp, paging, err := h.service.GetAdAccounts(fields, limit, after, before, false)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithPaging(c, "Successfully retrieved Meta ad accounts", resp, paging)
}
