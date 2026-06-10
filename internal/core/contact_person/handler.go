package contact_person

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/contact_person/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// FindAll godoc
// @Summary      List contact persons
// @Description  Get a paginated list of contact persons
// @Tags         Core - Contact Person
// @Accept       json
// @Produce      json
// @Param        page    query     int  false  "Page number" default(1)
// @Param        limit   query     int  false  "Items per page" default(10)
// @Success      200     {object}  response.PaginationResponse{data=[]dto.ContactPersonListResponse}
// @Failure      500     {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /contact-persons [get]
func (h *Handler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	res, total, err := h.service.FindAll(limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch contact persons", nil)
		return
	}

	meta := &response.PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}
	response.SuccessWithPagination(c, "Contact persons retrieved successfully", res, meta)
}

// FindByID godoc
// @Summary      Get contact person by ID
// @Description  Get contact person details by ID
// @Tags         Core - Contact Person
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Contact Person ID"
// @Success      200     {object}  response.SuccessResponse{data=dto.ContactPersonResponse}
// @Failure      404     {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /contact-persons/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	res, err := h.service.FindByID(id)
	if err != nil {
		if errors.Is(err, ErrContactPersonNotFound) {
			response.Error(c, http.StatusNotFound, "Contact person not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to fetch contact person", nil)
		return
	}

	response.Success(c, "Contact person retrieved successfully", res)
}

// Create godoc
// @Summary      Create contact person
// @Description  Create a new contact person
// @Tags         Core - Contact Person
// @Accept       json
// @Produce      json
// @Param        request body      dto.ContactPersonRequest  true  "Contact Person request"
// @Success      201     {object}  response.SuccessResponse{data=dto.ContactPersonResponse}
// @Failure      400     {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /contact-persons [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.ContactPersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create contact person", nil)
		return
	}

	response.Success(c, "Contact person created successfully", res)
}

// Update godoc
// @Summary      Update contact person
// @Description  Update an existing contact person
// @Tags         Core - Contact Person
// @Accept       json
// @Produce      json
// @Param        id      path      int                     true  "Contact Person ID"
// @Param        request body      dto.ContactPersonRequest  true  "Contact Person request"
// @Success      200     {object}  response.SuccessResponse{data=dto.ContactPersonResponse}
// @Failure      400     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /contact-persons/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	var req dto.ContactPersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	res, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, ErrContactPersonNotFound) {
			response.Error(c, http.StatusNotFound, "Contact person not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to update contact person", nil)
		return
	}

	response.Success(c, "Contact person updated successfully", res)
}

// Delete godoc
// @Summary      Delete contact person
// @Description  Delete a contact person by ID
// @Tags         Core - Contact Person
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Contact Person ID"
// @Success      200     {object}  response.SuccessResponse{data=nil}
// @Failure      404     {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /contact-persons/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, ErrContactPersonNotFound) {
			response.Error(c, http.StatusNotFound, "Contact person not found", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete contact person", nil)
		return
	}

	response.Success(c, "Contact person deleted successfully", nil)
}
