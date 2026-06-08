package adset

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.AdSetResponse{}

// GetAdSets godoc
// @Summary      Get AdSets
// @Description  Retrieve adsets from local database (synced from Meta)
// @Tags         Meta AdSets
// @Accept       json
// @Produce      json
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by adset name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.AdSetResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/adsets [get]
func (h *Handler) GetAdSets(c *gin.Context) {
	filter := AdSetFilter{
		CampaignID: c.Query("campaign_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAdSets(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved adsets", resp, meta)
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
