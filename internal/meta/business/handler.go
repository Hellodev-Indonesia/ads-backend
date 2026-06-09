package business

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetBusinesses godoc
// @Summary      Get Meta Businesses (Portfolios)
// @Description  Retrieve all business portfolios synced from Meta API
// @Tags         Meta Businesses
// @Accept       json
// @Produce      json
// @Param        search   query     string  false  "Search by name"
// @Param        page     query     int     false  "Page number"
// @Param        limit    query     int     false  "Items per page"
// @Success      200      {object}  response.Response{data=[]dto.BusinessResponse,meta=response.Meta}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /meta/businesses [get]
func (h *Handler) GetBusinesses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	search := c.Query("search")

	filter := BusinessFilter{
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	resp, meta, err := h.service.GetBusinesses(filter)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.SuccessWithMeta(c, "Successfully retrieved Meta businesses", resp, *meta)
}
