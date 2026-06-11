package campaign

import (
	"net/http"
	"strconv"

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

// GetCampaigns godoc
// @Summary      Get Campaigns
// @Description  Retrieve campaigns from local database (synced from Meta)
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Param        brand_id       query     int     false  "Filter by brand ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by campaign name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.CampaignResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/campaigns [get]
func (h *Handler) GetCampaigns(c *gin.Context) {
	adAccountID := c.Query("ad_account_id")

	var brandID *uint64
	if bid := c.Query("brand_id"); bid != "" {
		if id, err := strconv.ParseUint(bid, 10, 64); err == nil {
			brandID = &id
		}
	}

	filter := CampaignFilter{
		AccountID: adAccountID,
		BrandID:   brandID,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		Page:      parseQueryInt(c, "page", 1),
		Limit:     parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetCampaigns(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved campaigns", resp, meta)
}

func parseQueryInt(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(val)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}
