package response

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// SuccessResponse is used for Swagger documentation only
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
}

// ErrorResponse is used for Swagger documentation only
type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"Error message"`
	Errors  interface{} `json:"errors"`
}

// PaginationResponse is used for Swagger documentation only
type PaginationResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Total    int64 `json:"total"`
	LastPage int   `json:"last_page"`
}

func Success(c *gin.Context, message string, data interface{}) {
	if data == nil {
		data = gin.H{}
	} else {
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Slice && val.IsNil() {
			data = []interface{}{}
		}
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta Meta) {
	if data == nil {
		data = []interface{}{}
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(c *gin.Context, status int, message string, errors interface{}) {
	if errors == nil {
		errors = gin.H{}
	}
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, nil)
}
