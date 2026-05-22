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

// GetUnassignedAdAccounts godoc
// @Summary      Get Unassigned Meta Ad Accounts
// @Description  Retrieve all ad accounts not assigned to any brand
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Param        search   query     string  false  "Search by name"
// @Param        page     query     int     false  "Page number"
// @Param        limit    query     int     false  "Items per page"
// @Success      200      {object}  response.Response{data=[]dto.AdAccountResponse,meta=response.Meta}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/unassigned [get]
func (h *Handler) GetUnassignedAdAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	search := c.Query("search")

	filter := AdAccountFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	resp, meta, err := h.service.GetUnassigned(filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithMeta(c, "Successfully retrieved unassigned Meta ad accounts", resp, *meta)
}

// AssignBrand godoc
// @Summary      Assign Brand to Ad Account
// @Description  Assign a brand to a specific Meta ad account
// @Tags         Meta Ad Accounts
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Ad Account ID"
// @Param        request  body      dto.AssignBrandRequest  true  "Assign Brand Request"
// @Success      200      {object}  response.SuccessResponse
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/{id}/assign-brand [put]
func (h *Handler) AssignBrand(c *gin.Context) {
	id := c.Param("id")
	var req dto.AssignBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := h.service.AssignBrand(id, req.BrandID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully assigned brand to ad account", nil)
}

// UnassignBrand godoc
// @Summary      Unassign Brand from Ad Account
// @Description  Remove a brand assignment from a specific Meta ad account
// @Tags         Meta Ad Accounts
// @Produce      json
// @Param        id       path      string  true  "Ad Account ID"
// @Success      200      {object}  response.SuccessResponse
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/ad-accounts/{id}/unassign-brand [put]
func (h *Handler) UnassignBrand(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.UnassignBrand(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully unassigned brand from ad account", nil)
}
