package ads

import (
	"net/http"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.AdResponse{}
var _ = dto.CreativeResponse{}

func (h *Handler) getAdAccountID(c *gin.Context) string {
	id := c.Query("ad_account_id")
	if id != "" {
		return id
	}
	return config.MetaAdAccountID
}

// GetAds godoc
// @Summary      Get Ads
// @Description  Retrieve ads for the given or default ad account
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Param        fields         query     string  false  "Custom fields comma-separated preset"
// @Param        limit          query     string  false  "Pagination limit"
// @Param        after          query     string  false  "Cursor after"
// @Param        before         query     string  false  "Cursor before"
// @Success      200            {object}  response.Response{data=[]dto.AdResponse,paging=response.MetaPaging}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/ads [get]
func (h *Handler) GetAds(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	fields := c.Query("fields")
	if fields == "" {
		fields = DefaultAdFields
	}

	limit := c.Query("limit")
	after := c.Query("after")
	before := c.Query("before")

	resp, paging, err := h.service.GetAds(adAccountID, fields, limit, after, before, false)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithPaging(c, "Successfully retrieved ads", resp, paging)
}

// GetCreative godoc
// @Summary      Get Ad Creative
// @Description  Retrieve details of a specific ad creative
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        id       path      string  true   "Creative ID"
// @Param        fields   query     string  false  "Custom fields comma-separated preset"
// @Success      200      {object}  response.Response{data=dto.CreativeResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/creatives/{id} [get]
func (h *Handler) GetCreative(c *gin.Context) {
	creativeID := c.Param("id")
	if creativeID == "" {
		response.Error(c, http.StatusBadRequest, "Creative ID is required", nil)
		return
	}

	fields := c.Query("fields")
	if fields == "" {
		fields = DefaultCreativeFields
	}

	resp, err := h.service.GetCreative(creativeID, fields)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved ad creative", resp)
}
