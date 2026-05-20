package campaign

import (
	"net/http"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.CampaignResponse{}

func (h *Handler) getAdAccountID(c *gin.Context) string {
	id := c.Query("ad_account_id")
	if id != "" {
		return id
	}
	return config.MetaAdAccountID
}

// GetCampaigns godoc
// @Summary      Get Campaigns
// @Description  Retrieve campaigns for the given or default ad account
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Param        fields         query     string  false  "Custom fields comma-separated preset"
// @Param        limit          query     string  false  "Pagination limit"
// @Param        after          query     string  false  "Cursor after"
// @Param        before         query     string  false  "Cursor before"
// @Success      200            {object}  response.Response{data=[]dto.CampaignResponse,paging=response.MetaPaging}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/campaigns [get]
func (h *Handler) GetCampaigns(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	fields := c.Query("fields")
	if fields == "" {
		fields = DefaultFields
	}

	limit := c.Query("limit")
	after := c.Query("after")
	before := c.Query("before")

	resp, paging, err := h.service.GetCampaigns(adAccountID, fields, limit, after, before, false)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithPaging(c, "Successfully retrieved campaigns", resp, paging)
}
