package contact_person

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	contactPersons := router.Group("/contact-persons")
	contactPersons.Use(middleware.AuthMiddleware())
	{
		contactPersons.GET("", middleware.RequirePermission("core.contact_person.view"), handler.FindAll)
		contactPersons.GET("/:id", middleware.RequirePermission("core.contact_person.view"), handler.FindByID)
		contactPersons.POST("", middleware.RequirePermission("core.contact_person.create"), handler.Create)
		contactPersons.PUT("/:id", middleware.RequirePermission("core.contact_person.update"), handler.Update)
		contactPersons.DELETE("/:id", middleware.RequirePermission("core.contact_person.delete"), handler.Delete)
	}
}
