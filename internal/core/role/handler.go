package role

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/role/dto"
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
// @Summary      Get All Roles
// @Description  Retrieve all roles with permissions
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.SuccessResponse{data=[]dto.RoleResponse}
// @Router       /roles [get]
func (h *Handler) FindAll(c *gin.Context) {
	roles, err := h.service.FindAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, "Roles retrieved successfully", roles)
}

// FindByID godoc
// @Summary      Get Role by ID
// @Description  Retrieve a role by its ID
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Role ID"
// @Success      200      {object}  response.SuccessResponse{data=dto.RoleResponse}
// @Router       /roles/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	role, err := h.service.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Role not found", nil)
		return
	}
	response.Success(c, "Role retrieved successfully", role)
}

// Create godoc
// @Summary      Create Role
// @Description  Create a new role
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.RoleRequest  true  "Create Role Request"
// @Success      201      {object}  response.SuccessResponse
// @Router       /roles [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Role created successfully", nil)
}

// Update godoc
// @Summary      Update Role
// @Description  Update a role by ID
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Role ID"
// @Param        request  body      dto.RoleRequest  true  "Update Role Request"
// @Success      200      {object}  response.SuccessResponse
// @Router       /roles/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Role updated successfully", nil)
}

// Delete godoc
// @Summary      Delete Role
// @Description  Delete a role by ID
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Role ID"
// @Success      200      {object}  response.SuccessResponse
// @Router       /roles/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.Delete(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Role deleted successfully", nil)
}

// AssignPermissions godoc
// @Summary      Assign Permissions
// @Description  Assign permissions to a role
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "Role ID"
// @Param        request  body      dto.AssignPermissionRequest  true  "Assign Permission Request"
// @Success      200      {object}  response.SuccessResponse
// @Router       /roles/{id}/permissions [post]
func (h *Handler) AssignPermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, _ := strconv.Atoi(roleIDStr)

	var req dto.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.service.AssignPermissions(uint(roleID), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "Permissions assigned successfully", nil)
}
