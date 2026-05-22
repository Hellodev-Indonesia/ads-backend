package routes

import (
	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/internal/jobs"
	"github.com/alex/ads_backend/internal/meta/ad_account"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/dashboard"
	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/alex/ads_backend/internal/meta/sync"
	"github.com/alex/ads_backend/pkg/centrifugo"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/swagger"
	"github.com/alex/ads_backend/routes/core"
	"github.com/alex/ads_backend/routes/meta"
	"github.com/gin-gonic/gin"
)

func RegisterApiRoutes(router *gin.Engine) {
	// Documentation
	swagger.RegisterScalar(router, "Ads Backend API Reference", "/swagger.json")

	// Dev tool: Centrifugo listener page
	router.StaticFile("/centrifugo-listen", "./centrifugo-listen.html")

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

		// Register Core Routes
		core.RegisterAuthRoutes(v1, authHandler)
		core.RegisterUserRoutes(v1, userHandler)
		core.RegisterRoleRoutes(v1, roleHandler)
		core.RegisterPermissionRoutes(v1, permHandler)

		// --- META DOMAIN ---
		// Shared low-level client (single instance, injected into all sub-module services)
		metaClient := meta_client.NewClient(
			config.MetaGraphBaseURL,
			config.MetaGraphVersion,
			config.MetaAccessToken,
		)

		// Repositories (DB access layer)
		campaignRepo := campaign.NewRepository(config.DB)
		adSetRepo := adset.NewRepository(config.DB)
		adsRepo := ads.NewRepository(config.DB)
		insightRepo := insight.NewRepository(config.DB)
		syncRepo := sync.NewRepository(config.DB)

		// Services (Meta client + Repository)
		adAccountService := ad_account.NewService(metaClient)
		campaignService := campaign.NewService(metaClient, campaignRepo)
		adSetService := adset.NewService(metaClient, adSetRepo)
		adsService := ads.NewService(metaClient, adsRepo)
		insightService := insight.NewService(metaClient, insightRepo)
		syncService := sync.NewService(syncRepo)

		// Sub-module handlers
		adAccountHandler := ad_account.NewHandler(adAccountService)
		campaignHandler := campaign.NewHandler(campaignService)
		adSetHandler := adset.NewHandler(adSetService)
		adsHandler := ads.NewHandler(adsService)
		insightHandler := insight.NewHandler(insightService)

		dashboardRepo := dashboard.NewRepository(config.DB)
		dashboardService := dashboard.NewService(dashboardRepo)
		dashboardHandler := dashboard.NewHandler(dashboardService)

		// Register Meta Routes
		meta.RegisterMetaRoutes(v1, adAccountHandler, campaignHandler, adSetHandler, adsHandler, insightHandler, dashboardHandler)

		// Centrifugo publisher for real-time sync events
		centrifugoClient := centrifugo.NewClient(config.CentrifugoConfig.URL, config.CentrifugoConfig.APIKey)

		// Sync job (manual trigger only — no background ticker)
		syncJob := jobs.NewMetaAdsSyncJob(campaignService, adSetService, adsService, insightService, syncService, centrifugoClient)
		syncHandler := sync.NewHandler(syncJob, syncService)
		sync.RegisterRoutes(v1, syncHandler)
	}
}
