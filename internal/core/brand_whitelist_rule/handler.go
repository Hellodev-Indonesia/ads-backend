package brand_whitelist_rule

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create godoc
// @Summary      Create whitelist rule
// @Description  Create a new whitelist rule for a brand
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id path      int  true  "Brand ID"
// @Param        request  body      dto.CreateWhitelistRuleRequest  true  "Create Whitelist Rule Request"
// @Success      201      {object}  response.SuccessResponse{data=dto.WhitelistRuleResponse}
// @Router       /brands/{brand_id}/whitelist-rules [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	var req dto.CreateWhitelistRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rule, err := h.service.Create(brandID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Whitelist rule created", "data": rule})
}

// FindAll godoc
// @Summary      Get whitelist rules
// @Description  Get a list of whitelist rules for a brand
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id      path      int     true   "Brand ID"
// @Param        page          query     int     false  "Page number"
// @Param        limit         query     int     false  "Limit per page"
// @Param        scope         query     string  false  "Scope"
// @Param        match_type    query     string  false  "Match Type"
// @Param        is_active     query     bool    false  "Is Active"
// @Success      200  {object}  response.PaginationResponse{data=[]dto.WhitelistRuleResponse}
// @Router       /brands/{brand_id}/whitelist-rules [get]
// @Security     BearerAuth
func (h *Handler) FindAll(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	var filter dto.WhitelistRuleFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rules, total, err := h.service.FindAll(brandID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 25
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	lastPage := int(total) / limit
	if int(total)%limit > 0 {
		lastPage++
	}

	response.SuccessWithPagination(c, "Successfully retrieved whitelist rules", rules, &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	})
}

// FindByID godoc
// @Summary      Get whitelist rule by ID
// @Description  Get a whitelist rule by its ID
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id path      int  true  "Brand ID"
// @Param        id       path      int  true  "Whitelist Rule ID"
// @Success      200      {object}  response.SuccessResponse{data=dto.WhitelistRuleResponse}
// @Router       /brands/{brand_id}/whitelist-rules/{id} [get]
// @Security     BearerAuth
func (h *Handler) FindByID(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	rule, err := h.service.FindByID(brandID, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully retrieved whitelist rule", rule)
}

// Update godoc
// @Summary      Update whitelist rule
// @Description  Update an existing whitelist rule
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id path      int  true  "Brand ID"
// @Param        id       path      int  true  "Whitelist Rule ID"
// @Param        request  body      dto.UpdateWhitelistRuleRequest  true  "Update Whitelist Rule Request"
// @Success      200      {object}  response.SuccessResponse{data=dto.WhitelistRuleResponse}
// @Router       /brands/{brand_id}/whitelist-rules/{id} [put]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	var req dto.UpdateWhitelistRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rule, err := h.service.Update(brandID, id, req)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully updated whitelist rule", rule)
}

// Delete godoc
// @Summary      Delete whitelist rule
// @Description  Delete a whitelist rule
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id path      int  true  "Brand ID"
// @Param        id       path      int  true  "Whitelist Rule ID"
// @Success      200      {object}  response.SuccessResponse
// @Router       /brands/{brand_id}/whitelist-rules/{id} [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	if err := h.service.Delete(brandID, id); err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully deleted whitelist rule", nil)
}

// CheckURL godoc
// @Summary      Check URL against whitelist
// @Description  Check if a URL is allowed by the whitelist rules
// @Tags         Brand Whitelist Rules
// @Accept       json
// @Produce      json
// @Param        brand_id path      int  true  "Brand ID"
// @Param        request  body      dto.CheckURLRequest  true  "Check URL Request"
// @Success      200      {object}  response.SuccessResponse{data=dto.CheckURLResponse}
// @Router       /brands/{brand_id}/whitelist-rules/check-url [post]
// @Security     BearerAuth
func (h *Handler) CheckURL(c *gin.Context) {
	brandID, err := parseBrandID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand_id", nil)
		return
	}

	var req dto.CheckURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := h.service.IsURLAllowed(brandID, req.URL, req.Scope)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := dto.CheckURLResponse{Allowed: result.Allowed}
	if result.Rule != nil {
		r := toResponse(*result.Rule)
		resp.MatchedRule = &r
	}

	response.Success(c, "URL check completed", resp)
}

func parseBrandID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("brand_id"), 10, 64)
}

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
