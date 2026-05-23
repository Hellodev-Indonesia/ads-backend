# Graph Report - .  (2026-05-23)

## Corpus Check
- Corpus is ~40,052 words - fits in a single context window. You may not need a graph.

## Summary
- 785 nodes · 1088 edges · 88 communities (46 shown, 42 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 125 edges (avg confidence: 0.82)
- Token cost: 4,200 input · 1,800 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Meta Ad Creative & Delivery|Meta Ad Creative & Delivery]]
- [[_COMMUNITY_API Success Response Schemas|API Success Response Schemas]]
- [[_COMMUNITY_API Route Schema Docs|API Route Schema Docs]]
- [[_COMMUNITY_Ad Account Management|Ad Account Management]]
- [[_COMMUNITY_AdSet Sync & Formatting|AdSet Sync & Formatting]]
- [[_COMMUNITY_Campaign & AdSet Interfaces|Campaign & AdSet Interfaces]]
- [[_COMMUNITY_Brand DTOs|Brand DTOs]]
- [[_COMMUNITY_Brand & RBAC Service Layer|Brand & RBAC Service Layer]]
- [[_COMMUNITY_Infrastructure & Route Wiring|Infrastructure & Route Wiring]]
- [[_COMMUNITY_Meta Route Registrations|Meta Route Registrations]]
- [[_COMMUNITY_User & Permission Models|User & Permission Models]]
- [[_COMMUNITY_User & Role Handlers|User & Role Handlers]]
- [[_COMMUNITY_Sync Batch Service|Sync Batch Service]]
- [[_COMMUNITY_Insight Sync Service|Insight Sync Service]]
- [[_COMMUNITY_HTTP Response Utilities|HTTP Response Utilities]]
- [[_COMMUNITY_Auth & Brand API Docs|Auth & Brand API Docs]]
- [[_COMMUNITY_Sync Batch Repository|Sync Batch Repository]]
- [[_COMMUNITY_Campaign Sync Service|Campaign Sync Service]]
- [[_COMMUNITY_Dashboard Data Aggregation|Dashboard Data Aggregation]]
- [[_COMMUNITY_Brand Model & Concepts|Brand Model & Concepts]]
- [[_COMMUNITY_Auth Handler & DTOs|Auth Handler & DTOs]]
- [[_COMMUNITY_Ads Data Mapping|Ads Data Mapping]]
- [[_COMMUNITY_Realtime & Background Jobs|Realtime & Background Jobs]]
- [[_COMMUNITY_Entry Points & Bootstrap|Entry Points & Bootstrap]]
- [[_COMMUNITY_Permission Repository|Permission Repository]]
- [[_COMMUNITY_Role Repository|Role Repository]]
- [[_COMMUNITY_Swagger API Metadata|Swagger API Metadata]]
- [[_COMMUNITY_Database Config & Seeding|Database Config & Seeding]]
- [[_COMMUNITY_Meta API Client|Meta API Client]]
- [[_COMMUNITY_Brand Handler|Brand Handler]]
- [[_COMMUNITY_Dashboard Handler|Dashboard Handler]]
- [[_COMMUNITY_Sync Batch Handler|Sync Batch Handler]]
- [[_COMMUNITY_Role Service|Role Service]]
- [[_COMMUNITY_Campaign Repository|Campaign Repository]]
- [[_COMMUNITY_Insight Repository|Insight Repository]]
- [[_COMMUNITY_Permission Service|Permission Service]]
- [[_COMMUNITY_Permission Handler|Permission Handler]]
- [[_COMMUNITY_User Service|User Service]]
- [[_COMMUNITY_Auth Service & DTOs|Auth Service & DTOs]]
- [[_COMMUNITY_Config & JWT Utilities|Config & JWT Utilities]]
- [[_COMMUNITY_Brand List DTOs|Brand List DTOs]]
- [[_COMMUNITY_Campaign Handler|Campaign Handler]]
- [[_COMMUNITY_Sync Batch Model|Sync Batch Model]]
- [[_COMMUNITY_Insight Handler|Insight Handler]]
- [[_COMMUNITY_Ad Creative Models|Ad Creative Models]]
- [[_COMMUNITY_Architecture Documentation|Architecture Documentation]]
- [[_COMMUNITY_User Request DTOs|User Request DTOs]]
- [[_COMMUNITY_Role & Permission DTOs|Role & Permission DTOs]]
- [[_COMMUNITY_User DTOs|User DTOs]]
- [[_COMMUNITY_Centrifugo Client|Centrifugo Client]]
- [[_COMMUNITY_Auth DTOs|Auth DTOs]]
- [[_COMMUNITY_Permission DTOs|Permission DTOs]]
- [[_COMMUNITY_Dashboard Row DTOs|Dashboard Row DTOs]]
- [[_COMMUNITY_Upsert Pattern & Migrations|Upsert Pattern & Migrations]]
- [[_COMMUNITY_Sync Handler & Model|Sync Handler & Model]]
- [[_COMMUNITY_Ad Account Model|Ad Account Model]]
- [[_COMMUNITY_Meta Ad Model|Meta Ad Model]]
- [[_COMMUNITY_AdSet Model|AdSet Model]]
- [[_COMMUNITY_Alert Notification Model|Alert Notification Model]]
- [[_COMMUNITY_Brand Model|Brand Model]]
- [[_COMMUNITY_Whitelist Rule Model|Whitelist Rule Model]]
- [[_COMMUNITY_Campaign Model|Campaign Model]]
- [[_COMMUNITY_Fraud Log Model|Fraud Log Model]]
- [[_COMMUNITY_Insight Model|Insight Model]]
- [[_COMMUNITY_Ad Response DTOs|Ad Response DTOs]]
- [[_COMMUNITY_Whitelist Repository|Whitelist Repository]]
- [[_COMMUNITY_Ad Account Response DTO|Ad Account Response DTO]]
- [[_COMMUNITY_AdSet Response DTO|AdSet Response DTO]]
- [[_COMMUNITY_Assign Brand Request DTO|Assign Brand Request DTO]]
- [[_COMMUNITY_Campaign Response DTO|Campaign Response DTO]]
- [[_COMMUNITY_Insight Response DTO|Insight Response DTO]]
- [[_COMMUNITY_Sync Trigger Request DTO|Sync Trigger Request DTO]]
- [[_COMMUNITY_Creative Response DTO|Creative Response DTO]]
- [[_COMMUNITY_Sync Insight Request DTO|Sync Insight Request DTO]]
- [[_COMMUNITY_Permission Model|Permission Model]]
- [[_COMMUNITY_Role Model|Role Model]]
- [[_COMMUNITY_Alert & Whitelist Models|Alert & Whitelist Models]]
- [[_COMMUNITY_Whitelist Model|Whitelist Model]]
- [[_COMMUNITY_Ad Accounts Migration|Ad Accounts Migration]]
- [[_COMMUNITY_Brands Migration Rollback|Brands Migration Rollback]]
- [[_COMMUNITY_Auth Service Impl|Auth Service Impl]]
- [[_COMMUNITY_Sync Step Counts|Sync Step Counts]]
- [[_COMMUNITY_Response Wrapper|Response Wrapper]]

## God Nodes (most connected - your core abstractions)
1. `Success()` - 25 edges
2. `paths` - 23 edges
3. `responses` - 20 edges
4. `produces` - 17 edges
5. `tags` - 17 edges
6. `AuthMiddleware()` - 16 edges
7. `summary` - 16 edges
8. `parameters` - 15 edges
9. `RegisterApiRoutes()` - 14 edges
10. `SuccessWithPagination()` - 14 edges

## Surprising Connections (you probably didn't know these)
- `Asynq Background Job Worker` --semantically_similar_to--> `Centrifugo Channel 'meta:sync' (real-time sync notifications)`  [INFERRED] [semantically similar]
  cmd/worker/main.go → centrifugo-listen.html
- `Docker Compose (MySQL, Redis, Centrifugo services)` --conceptually_related_to--> `InitDB()`  [INFERRED]
  docker-compose.yml → config/database.go
- `Docker Compose (MySQL, Redis, Centrifugo services)` --conceptually_related_to--> `InitRedis()`  [INFERRED]
  docker-compose.yml → config/redis.go
- `Docker Compose (MySQL, Redis, Centrifugo services)` --conceptually_related_to--> `InitCentrifugo()`  [INFERRED]
  docker-compose.yml → config/centrifugo.go
- `RegisterApiRoutes()` --calls--> `NewMetaAdsSyncJob()`  [INFERRED]
  routes/api.go → internal/jobs/meta_ads_sync_job.go

## Hyperedges (group relationships)
- **API Server Bootstrap: LoadEnv → InitDB + InitRedis + InitMeta + InitCentrifugo → RegisterApiRoutes** — api_main_main, config_env_loadenv, config_database_initdb, config_redis_initredis, config_meta_initmeta, config_centrifugo_initcentrifugo [EXTRACTED 1.00]
- **Meta Sync Audit Pattern: meta_sync_batches → meta_sync_steps → meta Ads tables** — migration_003_sync_batches, migration_004_sync_steps, migration_002_meta_tables [INFERRED 0.85]
- **RBAC Schema: users + roles + permissions + user_roles + role_permissions** — migration_001_core_tables, claudemd_rbac, concept_soft_delete [EXTRACTED 1.00]
- **Brand as central FK anchor for ad_accounts, whitelist_rules, ad_creative_versions, fraud_logs, alerts** — migrations_000006_brands_table, migrations_000007_brand_to_ad_accounts, migrations_000008_brand_whitelist_rules_table, migrations_000010_ad_creative_versions_table, migrations_000011_fraud_logs_table, migrations_000011_alerts_table [EXTRACTED 1.00]
- **Auth login flow: Handler → Service → DTO** — auth_handler_login, auth_service_login, auth_dto_loginrequest, auth_dto_loginresponse [EXTRACTED 1.00]
- **Fraud detection flow: ad_creatives → ad_creative_versions → fraud_logs → alerts** — migrations_000009_ad_creatives_table, migrations_000010_ad_creative_versions_table, migrations_000011_fraud_logs_table, migrations_000011_alerts_table [INFERRED 0.85]
- **RBAC: User-Role-Permission many2many relationship** — user_model_user, role_model_role, permission_model_permission [EXTRACTED 1.00]
- **Brand CRUD layer: Repository-Service-DTO pattern** — brand_repository_repository, brand_service_service, brand_dto_brandfilter [EXTRACTED 0.95]
- **Fraud detection: BrandWhitelistRule + FraudLog linked by BrandID and MatchedRuleID** — brand_whitelist_rule_model_brandwhitelistrule, fraud_log_model_fraudlog, brand_repository_repository [INFERRED 0.75]
- **Meta Ads Full Sync Pipeline** — jobs_metaadssynctask, ad_account_service, ads_service, adset_repository [EXTRACTED 0.95]
- **AdAccount CRUD Pattern (Handler-Service-Repository)** — ad_account_handler, ad_account_service, ad_account_repository [EXTRACTED 1.00]
- **Ads CRUD Pattern (Handler-Service-Repository)** — ads_handler, ads_service, ads_repository [EXTRACTED 1.00]
- **Campaign Read Flow: Handler calls Service calls Repository** — campaign_handler_handler, campaign_service_serviceimpl, campaign_repository_repository [EXTRACTED 1.00]
- **Insight Sync Flow: ServiceImpl fetches from Meta API and upserts via Repository** — insight_service_serviceimpl, insight_repository_repository, insight_dto_syncrequest_syncinsightrequest [EXTRACTED 1.00]
- **Dashboard aggregates Campaign, AdSet and Insight data via SQL JOINs in Repository** — dashboard_repository_campaigndashboardscan, dashboard_repository_adsetdashboardscan, dashboard_repository_addashboardscan [INFERRED 0.85]
- **JWT Auth + Permission Check Flow** — middleware_auth_authmiddleware, middleware_auth_requirepermission, utils_jwt_validatetoken, utils_jwt_jwtclaims, response_response_unauthorized, response_response_forbidden [EXTRACTED 1.00]
- **Meta Domain Dependency Injection in RegisterApiRoutes** — routes_api_registerapiroutes, meta_client_client_client, sync_service_service, centrifugo_client_client [EXTRACTED 1.00]
- **Sync Batch Lifecycle (Repository, Service, Route)** — sync_repository_repository, sync_service_service, sync_route_registerroutes, sync_dto_triggersyncrequest [EXTRACTED 1.00]

## Communities (88 total, 42 thin omitted)

### Community 0 - "Meta Ad Creative & Delivery"
Cohesion: 0.05
Nodes (25): AdCreative Model, AdCreativeVersion Model, AdFilter, AdResponse DTO, CreativeRef DTO, CreativeResponse DTO, Handler, parseQueryInt() (+17 more)

### Community 1 - "API Success Response Schemas"
Cohesion: 0.06
Nodes (47): description, schema, description, schema, description, schema, description, schema (+39 more)

### Community 2 - "API Route Schema Docs"
Cohesion: 0.16
Nodes (44): description, get, get, consumes, description, parameters, produces, responses (+36 more)

### Community 3 - "Ad Account Management"
Cohesion: 0.07
Nodes (10): AdAccountFilter, AssignBrandRequest DTO, AdAccountResponse DTO, Handler, MetaAdAccount Model, Repository, AdAccount Routes, Service (+2 more)

### Community 4 - "AdSet Sync & Formatting"
Cohesion: 0.13
Nodes (20): Service, formatDecimal(), formatTime(), mapDTOToModel(), mapModelToDTO(), parseDecimal(), parseTime(), serviceImpl (+12 more)

### Community 5 - "Campaign & AdSet Interfaces"
Cohesion: 0.10
Nodes (29): AdSetResponse DTO, AdSet Service Interface, AdSet serviceImpl, CampaignResponse DTO, Campaign Handler, MetaCampaign Model, CampaignFilter, Campaign Repository Interface (+21 more)

### Community 6 - "Brand DTOs"
Cohesion: 0.10
Nodes (13): dto.BrandFilter, dto.BrandResponse, dto.CreateBrandRequest, dto.UpdateBrandRequest, Repository, FilterBrand(), brand.Repository (interface), brand.repository (struct) (+5 more)

### Community 7 - "Brand & RBAC Service Layer"
Cohesion: 0.13
Nodes (20): brand.Service (interface), dto.PermissionFilter, dto.PermissionRequest, dto.PermissionResponse, permission.Handler, permission.Repository (interface), permission.repository (struct), RegisterRoutes() (+12 more)

### Community 8 - "Infrastructure & Route Wiring"
Cohesion: 0.11
Nodes (16): centrifugo.Client, Manual Dependency Injection Pattern, Meta Graph API Pagination, RegisterCoreRoutes(), meta_client.BaseResponse, meta_client.Client, meta_client.Error, RegisterMetaRoutes() (+8 more)

### Community 9 - "Meta Route Registrations"
Cohesion: 0.14
Nodes (11): RegisterRoutes(), RegisterRoutes(), RegisterRoutes(), RegisterRoutes(), RegisterRoutes(), RBAC Permission Check Pattern, RegisterRoutes(), RegisterRoutes() (+3 more)

### Community 10 - "User & Permission Models"
Cohesion: 0.10
Nodes (10): Permission, Role, user.Handler, FilterUser(), User, Repository, user.Repository (interface), user.repository (struct) (+2 more)

### Community 11 - "User & Role Handlers"
Cohesion: 0.17
Nodes (3): Success(), Handler, Handler

### Community 12 - "Sync Batch Service"
Cohesion: 0.17
Nodes (5): Service, calculateDurationMs(), generateBatchCode(), nullableString(), StartBatchInput

### Community 13 - "Insight Sync Service"
Cohesion: 0.22
Nodes (7): Service, formatDecimal(), mapDTOToModel(), mapModelToDTO(), parseDecimal(), parseInt(), serviceImpl

### Community 14 - "HTTP Response Utilities"
Cohesion: 0.19
Nodes (13): ErrorResponse, Meta, MetaPaging, PaginationMeta, PaginationResponse, Response, BadRequest(), Error() (+5 more)

### Community 15 - "Auth & Brand API Docs"
Cohesion: 0.30
Nodes (14): post, post, post, post, /auth/login, /auth/logout, /meta/sync, consumes (+6 more)

### Community 17 - "Campaign Sync Service"
Cohesion: 0.24
Nodes (8): Service, formatDecimal(), formatTime(), mapDTOToModel(), mapModelToDTO(), parseDecimal(), parseTime(), serviceImpl

### Community 18 - "Dashboard Data Aggregation"
Cohesion: 0.26
Nodes (7): adDashboardScan, adSetDashboardScan, campaignDashboardScan, DashboardFilter, Repository, intToStr(), javaStrToInt()

### Community 19 - "Brand Model & Concepts"
Cohesion: 0.25
Nodes (11): brand.Brand (model), Ad creative versioning with change tracking, Brand as multi-tenancy boundary, Fraud detection pipeline (logs + alerts), brands table (migration 000006), Add brand_id to meta_ad_accounts (migration 000007), brand_whitelist_rules table (migration 000008), ad_creatives table (migration 000009) (+3 more)

### Community 20 - "Auth Handler & DTOs"
Cohesion: 0.22
Nodes (6): dto.LoginRequest, Handler, auth.Handler, RegisterRoutes(), auth.Service (interface), Swagger/OpenAPI docs (generated)

### Community 21 - "Ads Data Mapping"
Cohesion: 0.29
Nodes (5): formatTime(), mapDTOToModel(), mapModelToDTO(), parseTime(), serviceImpl

### Community 22 - "Realtime & Background Jobs"
Cohesion: 0.20
Nodes (8): Centrifugo WebSocket Listener (test/debug HTML page), Asynq Background Job Worker, Centrifugo Channel 'meta:sync' (real-time sync notifications), CentrifugoConfig (struct var), InitCentrifugo(), InitRedis(), RedisClient (global *redis.Client), Docker Compose (MySQL, Redis, Centrifugo services)

### Community 23 - "Entry Points & Bootstrap"
Cohesion: 0.22
Nodes (6): main(), Meta Graph API v25.0, LoadEnv(), InitMeta(), MetaGraphBaseURL / MetaAccessToken / MetaAdAccountID (globals), main()

### Community 26 - "Swagger API Metadata"
Cohesion: 0.22
Nodes (8): basePath, host, info, contact, description, title, version, swagger

### Community 27 - "Database Config & Seeding"
Cohesion: 0.25
Nodes (5): DB (global *gorm.DB), InitDB(), main(), SeedCore(), Run()

### Community 28 - "Meta API Client"
Cohesion: 0.22
Nodes (5): BaseResponse, Client, Error, errorWrapper, Paging

### Community 30 - "Dashboard Handler"
Cohesion: 0.36
Nodes (3): Handler, parseQueryInt(), SuccessWithPagination()

### Community 31 - "Sync Batch Handler"
Cohesion: 0.25
Nodes (3): Handler, parseIntQuery(), JobTrigger

### Community 38 - "Auth Service & DTOs"
Cohesion: 0.29
Nodes (4): dto.AuthUserResponse, dto.LoginResponse, Service, Super Admin full permission bypass pattern

### Community 39 - "Config & JWT Utilities"
Cohesion: 0.43
Nodes (6): GetEnv(), GenerateCentrifugoToken(), GenerateToken(), utils.JWTClaims, ValidateToken(), JWTClaims

### Community 40 - "Brand List DTOs"
Cohesion: 0.33
Nodes (5): BrandFilter, BrandListResponse, BrandResponse, CreateBrandRequest, UpdateBrandRequest

### Community 42 - "Sync Batch Model"
Cohesion: 0.33
Nodes (3): MetaSyncBatch, MetaSyncStep, StepCounts

### Community 45 - "Architecture Documentation"
Cohesion: 0.40
Nodes (5): Go REST API Architecture (Meta Ads Backend), Repository → Service → Handler Layer Pattern, RBAC with domain.module.action Permission Naming, Soft Delete Pattern (deleted_at on core tables), Migration 000001: Core Tables (users, roles, permissions, user_roles, role_permissions)

### Community 46 - "User Request DTOs"
Cohesion: 0.40
Nodes (4): RoleBrief, UserFilter, UserRequest, UserResponse

### Community 47 - "Role & Permission DTOs"
Cohesion: 0.40
Nodes (4): AssignPermissionRequest, RoleFilter, RoleRequest, RoleResponse

### Community 48 - "User DTOs"
Cohesion: 0.40
Nodes (5): RoleBrief DTO, UserFilter DTO, UserRequest DTO, UserResponse DTO, User Service

### Community 50 - "Auth DTOs"
Cohesion: 0.50
Nodes (3): AuthUserResponse, LoginRequest, LoginResponse

### Community 51 - "Permission DTOs"
Cohesion: 0.50
Nodes (3): PermissionFilter, PermissionRequest, PermissionResponse

### Community 52 - "Dashboard Row DTOs"
Cohesion: 0.50
Nodes (3): AdDashboardRow, AdSetDashboardRow, CampaignDashboardRow

### Community 53 - "Upsert Pattern & Migrations"
Cohesion: 0.50
Nodes (4): Upsert Pattern via Unique Constraint on meta_insights, Migration 000002: Meta Tables (meta_campaigns, meta_ad_sets, meta_ads, meta_insights), Migration 000003: meta_sync_batches Table, Migration 000004: meta_sync_steps Table

### Community 54 - "Sync Handler & Model"
Cohesion: 0.50
Nodes (4): Sync Handler, JobTrigger Interface, MetaSyncBatch Model, MetaSyncStep Model

## Knowledge Gaps
- **132 isolated node(s):** `Role`, `RoleRequest`, `RoleResponse`, `AssignPermissionRequest`, `RoleFilter` (+127 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `AuthMiddleware()` connect `Meta Route Registrations` to `Brand & RBAC Service Layer`, `Config & JWT Utilities`, `User & Permission Models`, `HTTP Response Utilities`, `Auth Handler & DTOs`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Why does `Success()` connect `User & Role Handlers` to `Meta Ad Creative & Delivery`, `Ad Account Management`, `Permission Handler`, `HTTP Response Utilities`, `Auth Handler & DTOs`, `Brand Handler`, `Sync Batch Handler`?**
  _High betweenness centrality (0.059) - this node is a cross-community bridge._
- **Why does `SuccessWithPagination()` connect `Dashboard Handler` to `Meta Ad Creative & Delivery`, `Permission Handler`, `Campaign Handler`, `User & Role Handlers`, `Insight Handler`, `HTTP Response Utilities`, `Brand Handler`, `Sync Batch Handler`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `Success()` (e.g. with `.FindByID()` and `.Create()`) actually correct?**
  _`Success()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Role`, `RoleRequest`, `RoleResponse` to the rest of the system?**
  _138 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Meta Ad Creative & Delivery` be split into smaller, more focused modules?**
  _Cohesion score 0.054693877551020405 - nodes in this community are weakly interconnected._
- **Should `API Success Response Schemas` be split into smaller, more focused modules?**
  _Cohesion score 0.06290471785383904 - nodes in this community are weakly interconnected._