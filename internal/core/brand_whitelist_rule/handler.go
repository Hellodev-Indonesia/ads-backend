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
