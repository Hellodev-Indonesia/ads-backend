package core

import (
	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/internal/core/brand"
	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule"
	fraudlog "github.com/alex/ads_backend/internal/core/fraud_log"
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/internal/notification/alert"
	"github.com/gin-gonic/gin"
)

func RegisterCoreRoutes(
	r *gin.RouterGroup,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	roleHandler *role.Handler,
	permHandler *permission.Handler,
	brandHandler *brand.Handler,
	whitelistHandler *brand_whitelist_rule.Handler,
	fraudLogHandler *fraudlog.Handler,
	alertHandler *alert.Handler,
) {
	auth.RegisterRoutes(r, authHandler)
	user.RegisterRoutes(r, userHandler)
	role.RegisterRoutes(r, roleHandler)
	permission.RegisterRoutes(r, permHandler)
	brand.RegisterRoutes(r, brandHandler)
	brand_whitelist_rule.RegisterRoutes(r, whitelistHandler)
	fraudlog.RegisterRoutes(r, fraudLogHandler)
	alert.RegisterRoutes(r, alertHandler)
}
