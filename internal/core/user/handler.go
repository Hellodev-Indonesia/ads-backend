package user

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/user/dto"
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
// @Summary      Register User
// @Description  Create a new user
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.UserRequest  true  "Create User Request"
// @Success      201      {object}  response.SuccessResponse
// @Router       /users [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "User created successfully", nil)
}

// Update godoc
// @Summary      Update User
// @Description  Update a user by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "User ID"
// @Param        request  body      dto.UserRequest  true  "Update User Request"
// @Success      200      {object}  response.SuccessResponse
// @Router       /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "User updated successfully", nil)
}

// Delete godoc
// @Summary      Delete User
// @Description  Delete a user by ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "User ID"
// @Success      200      {object}  response.SuccessResponse
// @Router       /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.Delete(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// FindAll godoc
// @Summary      Get All Users
// @Description  Retrieve all users
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name     query     string  false  "Filter by name"
// @Param        email    query     string  false  "Filter by email"
// @Param        page     query     int     false  "Page number" default(1)
// @Param        limit    query     int     false  "Items per page" default(25)
// @Success      200      {object}  response.PaginationResponse{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *Handler) FindAll(c *gin.Context) {
	var filter dto.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	users, meta, err := h.service.FindAll(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Users retrieved successfully", users, meta)
}

// FindByID godoc
// @Summary      Get User by ID
// @Description  Retrieve a user by its ID
// @Tags         User Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "User ID"
// @Success      200      {object}  response.SuccessResponse{data=dto.UserResponse}
// @Router       /users/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, err := h.service.FindByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", nil)
		return
	}
	response.Success(c, "User retrieved successfully", user)
}
