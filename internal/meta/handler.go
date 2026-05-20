package meta

import (
	"net/http"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

var _ = dto.AdAccountResponse{}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) getAdAccountID(c *gin.Context) string {
	id := c.Query("ad_account_id")
	if id != "" {
		return id
	}
	return config.MetaAdAccountID
}

// GetAdAccounts godoc
// @Summary      Get Meta Ad Accounts
// @Description  Retrieve all ad accounts associated with the system user token
// @Tags         Meta Marketing
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

// GetCampaigns godoc
// @Summary      Get Campaigns
// @Description  Retrieve campaigns for the given or default ad account
// @Tags         Meta Marketing
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Success      200            {object}  response.SuccessResponse{data=[]dto.CampaignResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/campaigns [get]
func (h *Handler) GetCampaigns(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	resp, err := h.service.GetCampaigns(adAccountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved campaigns", resp)
}

// GetAdSets godoc
// @Summary      Get AdSets
// @Description  Retrieve adsets for the given or default ad account
// @Tags         Meta Marketing
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Success      200            {object}  response.SuccessResponse{data=[]dto.AdSetResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/adsets [get]
func (h *Handler) GetAdSets(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	resp, err := h.service.GetAdSets(adAccountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved adsets", resp)
}

// GetAds godoc
// @Summary      Get Ads
// @Description  Retrieve ads for the given or default ad account
// @Tags         Meta Marketing
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Success      200            {object}  response.SuccessResponse{data=[]dto.AdResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/ads [get]
func (h *Handler) GetAds(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	resp, err := h.service.GetAds(adAccountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved ads", resp)
}

// GetInsights godoc
// @Summary      Get Insights
// @Description  Retrieve today's insights for the given or default ad account
// @Tags         Meta Marketing
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Success      200            {object}  response.SuccessResponse{data=[]dto.InsightResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/insights [get]
func (h *Handler) GetInsights(c *gin.Context) {
	adAccountID := h.getAdAccountID(c)
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	resp, err := h.service.GetInsights(adAccountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved insights", resp)
}
