package brand

import (
	"net/http"

	"github.com/alex/ads_backend/internal/core/brand/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// FindAll godoc
// @Summary      List Brands
// @Description  Get a paginated list of brands
// @Tags         Brand
// @Produce      json
// @Param        page      query     int     false  "Page number" default(1)
// @Param        limit     query     int     false  "Items per page" default(25)
// @Param        name      query     string  false  "Search by name"
// @Param        is_active query     bool    false  "Filter by active status"
// @Success      200       {object}  response.PaginationResponse{data=[]dto.BrandResponse}
// @Failure      500       {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /core/brands [get]
func (h *Handler) FindAll(c *gin.Context) {
	var filter dto.BrandFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brands, total, err := h.service.FindAll(filter)
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

	response.SuccessWithPagination(c, "Successfully retrieved brands", brands, &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	})
}

// FindBySlug godoc
// @Summary      Get Brand Details
// @Description  Get details of a specific brand by Slug
// @Tags         Brand
// @Produce      json
// @Param        slug   path      string  true  "Brand Slug"
// @Success      200  {object}  response.SuccessResponse{data=dto.BrandResponse}
// @Failure      404  {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /core/brands/{slug} [get]
func (h *Handler) FindBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Error(c, http.StatusBadRequest, "Invalid slug", nil)
		return
	}

	brand, err := h.service.FindBySlug(slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully retrieved brand", brand)
}

// Create godoc
// @Summary      Create Brand
// @Description  Create a new brand with optional photo upload
// @Tags         Brand
// @Accept       multipart/form-data
// @Produce      json
// @Param        name         formData  string  true  "Brand Name"
// @Param        photo        formData  file    false "Brand Photo (max 10MB)"
// @Param        description  formData  string  false "Brand Description"
// @Param        is_active    formData  bool    false "Is Active"
// @Success      201      {object}  response.SuccessResponse{data=dto.BrandResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /core/brands [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brand, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully created brand", brand)
}

// Update godoc
// @Summary      Update Brand
// @Description  Update an existing brand with optional photo upload
// @Tags         Brand
// @Accept       multipart/form-data
// @Produce      json
// @Param        slug         path      string  true  "Brand Slug"
// @Param        name         formData  string  false "Brand Name"
// @Param        photo        formData  file    false "Brand Photo (max 10MB)"
// @Param        description  formData  string  false "Brand Description"
// @Param        is_active    formData  bool    false "Is Active"
// @Success      200      {object}  response.SuccessResponse{data=dto.BrandResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /core/brands/{slug} [put]
func (h *Handler) Update(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Error(c, http.StatusBadRequest, "Invalid slug", nil)
		return
	}

	var req dto.UpdateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	brand, err := h.service.Update(slug, req)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully updated brand", brand)
}

// Delete godoc
// @Summary      Delete Brand
// @Description  Soft delete a brand
// @Tags         Brand
// @Produce      json
// @Param        slug   path      string  true  "Brand Slug"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /core/brands/{slug} [delete]
func (h *Handler) Delete(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Error(c, http.StatusBadRequest, "Invalid slug", nil)
		return
	}

	if err := h.service.Delete(slug); err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully deleted brand", nil)
}
