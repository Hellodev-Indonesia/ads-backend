package routes

import (
	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/internal/jobs"
	metaInternal "github.com/alex/ads_backend/internal/meta"
	"github.com/alex/ads_backend/pkg/swagger"
	"github.com/alex/ads_backend/routes/core"
	metaRoutes "github.com/alex/ads_backend/routes/meta"
	"github.com/gin-gonic/gin"
)

func RegisterApiRoutes(router *gin.Engine) {
	// Documentation
	swagger.RegisterScalar(router, "Ads Backend API Reference", "/swagger.json")

	v1 := router.Group("/api/v1")
	{
		// --- CORE DOMAIN ---
		// Repositories
		userRepo := user.NewRepository(config.DB)
		roleRepo := role.NewRepository(config.DB)
		permRepo := permission.NewRepository(config.DB)

		// Services
		permService := permission.NewService(permRepo)
		roleService := role.NewService(roleRepo, permRepo)
		userService := user.NewService(userRepo, roleRepo)
		authService := auth.NewService(userRepo, permRepo)

		// Handlers
		authHandler := auth.NewHandler(authService)
		userHandler := user.NewHandler(userService)
		roleHandler := role.NewHandler(roleService)
		permHandler := permission.NewHandler(permService)

		// Register Modular Routes
		core.RegisterAuthRoutes(v1, authHandler)
		core.RegisterUserRoutes(v1, userHandler)
		core.RegisterRoleRoutes(v1, roleHandler)
		core.RegisterPermissionRoutes(v1, permHandler)

		// --- META DOMAIN ---
		metaClient := metaInternal.NewClient()
		metaService := metaInternal.NewService(metaClient)
		metaHandler := metaInternal.NewHandler(metaService)

		// Register Meta Routes
		metaRoutes.RegisterMetaRoutes(v1, metaHandler)

		// Start Meta Background Sync Job
		syncJob := jobs.NewMetaAdsSyncJob(metaService)
		syncJob.Start()
	}
}
