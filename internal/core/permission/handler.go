package permission

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/permission/dto"
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
// @Summary      Get All Permissions
// @Description  Retrieve all permissions
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name   query     string  false  "Filter by name"
// @Param        page   query     int     false  "Page number" default(1)
// @Param        limit  query     int     false  "Items per page" default(25)
// @Success      200    {object}  response.PaginationResponse{data=[]dto.PermissionResponse}
// @Router       /permissions [get]
func (h *Handler) FindAll(c *gin.Context) {
	var filter dto.PermissionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	permissions, meta, err := h.service.FindAll(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Permissions retrieved successfully", permissions, meta)
}

// FindByID godoc
// @Summary      Get Permission by ID
// @Description  Retrieve a permission by its ID
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Permission ID"
// @Success      200      {object}  response.SuccessResponse{data=dto.PermissionResponse}
// @Router       /permissions/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	p, err := h.service.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Permission not found", nil)
		return
	}
	response.Success(c, "Permission retrieved successfully", p)
}

// Create godoc
// @Summary      Create Permission
// @Description  Create a new permission
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.PermissionRequest  true  "Create Permission Request"
// @Success      201      {object}  response.SuccessResponse
// @Router       /permissions [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Permission created successfully", nil)
}

// Update godoc
// @Summary      Update Permission
// @Description  Update a permission by ID
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Permission ID"
// @Param        request  body      dto.PermissionRequest  true  "Update Permission Request"
// @Success      200      {object}  response.SuccessResponse
// @Router       /permissions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Permission updated successfully", nil)
}

// Delete godoc
// @Summary      Delete Permission
// @Description  Delete a permission by ID
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Permission ID"
// @Success      200      {object}  response.SuccessResponse
// @Router       /permissions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.Delete(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Permission deleted successfully", nil)
}
