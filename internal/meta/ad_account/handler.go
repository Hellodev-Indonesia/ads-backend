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
// @Param        brand_id query     int     false  "Filter by Brand ID"
// @Param        business_id query  string  false  "Filter by Business ID"
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
	
	var brandID *uint64
	brandIDStr := c.Query("brand_id")
	if brandIDStr != "" {
		if id, err := strconv.ParseUint(brandIDStr, 10, 64); err == nil {
			brandID = &id
		}
	}

	businessID := c.Query("business_id")
	var businessIDPtr *string
	if businessID != "" {
		businessIDPtr = &businessID
	}

	filter := AdAccountFilter{
		Search:     search,
		Page:       page,
		Limit:      limit,
		BrandID:    brandID,
		BusinessID: businessIDPtr,
	}

	resp, meta, err := h.service.GetAdAccounts(filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithMeta(c, "Successfully retrieved Meta ad accounts", resp, *meta)
}

// GetUnassignedAdAccounts godoc
// @Summary      Get Unassigned Ad Accounts
// @Description  Get paginated list of Meta ad accounts that are not assigned to any brand
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Param        search    query     string  false  "Search by name"
// @Param        page      query     int     false  "Page number"
// @Param        limit     query     int     false  "Items per page"
// @Param        business_id query   string  false  "Filter by Business ID"
// @Success      200       {object}  response.PaginationResponse{data=[]dto.AdAccountResponse}
// @Failure      400       {object}  response.ErrorResponse
// @Failure      401       {object}  response.ErrorResponse
// @Failure      500       {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/unassigned [get]
func (h *Handler) GetUnassignedAdAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	businessID := c.Query("business_id")
	var businessIDPtr *string
	if businessID != "" {
		businessIDPtr = &businessID
	}

	filter := AdAccountFilter{
		Search:     c.Query("search"),
		Page:       page,
		Limit:      limit,
		BusinessID: businessIDPtr,
	}

	resp, meta, err := h.service.GetUnassigned(filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithMeta(c, "Successfully retrieved unassigned Meta ad accounts", resp, *meta)
}

// AssignBrand godoc
// @Summary      Assign or Unassign Brand to Ad Accounts
// @Description  Assign a brand to multiple Meta ad accounts (or unassign if brand_id is null)
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AssignBrandRequest  true  "Assign Brand Request"
// @Success      200      {object}  response.SuccessResponse
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/brand [put]
func (h *Handler) AssignBrand(c *gin.Context) {
	var req dto.AssignBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if len(req.AdAccountIDs) == 0 {
		response.Error(c, http.StatusBadRequest, "ad_account_ids is required", nil)
		return
	}

	if err := h.service.BulkAssignBrand(req.AdAccountIDs, req.BrandID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if req.BrandID == nil {
		response.Success(c, "Successfully unassigned brand from ad accounts", nil)
		return
	}

	response.Success(c, "Successfully assigned brand to ad accounts", nil)
}

// DisconnectBrand godoc
// @Summary      Disconnect Brand from Ad Account
// @Description  Disconnect a brand from a specific Meta ad account
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Ad Account ID"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/{id}/disconnect [put]
func (h *Handler) DisconnectBrand(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "ad account id is required", nil)
		return
	}

	if err := h.service.DisconnectBrand(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully disconnected brand from ad account", nil)
}

// GetBusinessOptions gets the unique business options
// @Summary Get Business Options
// @Description Get unique business options available in ad accounts
// @Tags Meta Ad Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]dto.BusinessOptionResponse}
// @Router /meta/ad-accounts/businesses [get]
func (h *Handler) GetBusinessOptions(c *gin.Context) {
	businesses, err := h.service.GetBusinessOptions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, "Business options retrieved successfully", businesses)
}
