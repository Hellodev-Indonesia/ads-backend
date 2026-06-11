import re

def read_file(path):
    with open(path, 'r') as f: return f.read()

def write_file(path, content):
    with open(path, 'w') as f: f.write(content)

# 1. ad_account repository
adacc_repo = read_file('internal/meta/ad_account/repository.go')
brand_scan = re.search(r'(type brandDashboardScan struct \{.*?\n\})', adacc_repo, re.MULTILINE | re.DOTALL)
brand_repo_func = re.search(r'(func \(r \*repository\) FindBrandDashboard.*?^})', adacc_repo, re.MULTILINE | re.DOTALL)

if brand_scan and brand_repo_func:
    adacc_repo = adacc_repo.replace(brand_scan.group(1), '')
    adacc_repo = adacc_repo.replace(brand_repo_func.group(1), '')
    adacc_repo = adacc_repo.replace('\tFindBrandDashboard(filter AdAccountFilter) ([]brandDashboardScan, int64, error)\n', '')
    write_file('internal/meta/ad_account/repository.go', adacc_repo)

    # Move to core/brand
    brand_repo = read_file('internal/core/brand/repository.go')
    brand_repo = brand_repo.replace('FindAll(filter dto.BrandFilter) ([]Brand, int64, error)\n}', 'FindAll(filter dto.BrandFilter) ([]Brand, int64, error)\n\tFindBrandDashboard(filter dto.BrandDashboardFilter) ([]brandDashboardScan, int64, error)\n}')
    
    brand_repo_func_text = brand_repo_func.group(1).replace('AdAccountFilter', 'dto.BrandDashboardFilter')
    brand_repo += '\n' + brand_scan.group(1) + '\n' + brand_repo_func_text + '\n'
    write_file('internal/core/brand/repository.go', brand_repo)

# 2. ad_account service
adacc_svc = read_file('internal/meta/ad_account/service.go')
brand_svc_func = re.search(r'(func \(s \*serviceImpl\) GetBrandDashboard.*?^})', adacc_svc, re.MULTILINE | re.DOTALL)

if brand_svc_func:
    adacc_svc = adacc_svc.replace(brand_svc_func.group(1), '')
    adacc_svc = adacc_svc.replace('\tGetBrandDashboard(filter AdAccountFilter) ([]dto.BrandDashboardResponse, *response.PaginationMeta, error)\n', '')
    write_file('internal/meta/ad_account/service.go', adacc_svc)

    # Move to core/brand
    brand_svc = read_file('internal/core/brand/service.go')
    brand_svc = brand_svc.replace('FindAll(filter dto.BrandFilter) ([]dto.BrandResponse, *response.PaginationMeta, error)\n}', 'FindAll(filter dto.BrandFilter) ([]dto.BrandResponse, *response.PaginationMeta, error)\n\tGetBrandDashboard(filter dto.BrandDashboardFilter) ([]dto.BrandDashboardResponse, *response.PaginationMeta, error)\n}')
    
    brand_svc_func_text = brand_svc_func.group(1).replace('AdAccountFilter', 'dto.BrandDashboardFilter')
    brand_svc += '\n' + brand_svc_func_text + '\n'
    write_file('internal/core/brand/service.go', brand_svc)

# 3. ad_account handler
adacc_hndlr = read_file('internal/meta/ad_account/handler.go')
brand_hndlr_func = re.search(r'(// GetBrandDashboard godoc.*?^})', adacc_hndlr, re.MULTILINE | re.DOTALL)
parseQueryInt = re.search(r'(func parseQueryInt.*?^})', adacc_hndlr, re.MULTILINE | re.DOTALL)

if brand_hndlr_func:
    adacc_hndlr = adacc_hndlr.replace(brand_hndlr_func.group(1), '')
    if parseQueryInt:
        adacc_hndlr = adacc_hndlr.replace(parseQueryInt.group(1), '')
    write_file('internal/meta/ad_account/handler.go', adacc_hndlr)

    # Move to core/brand
    brand_hndlr = read_file('internal/core/brand/handler.go')
    
    brand_hndlr_func_text = brand_hndlr_func.group(1).replace('AdAccountFilter', 'dto.BrandDashboardFilter').replace('// @Router       /meta/brands [get]', '// @Router       /core/brands/dashboard [get]')
    brand_hndlr += '\n' + brand_hndlr_func_text + '\n'
    if parseQueryInt:
        brand_hndlr += '\n' + parseQueryInt.group(1) + '\n'
    write_file('internal/core/brand/handler.go', brand_hndlr)

# 4. ad_account route
adacc_route = read_file('internal/meta/ad_account/route.go')
adacc_route = adacc_route.replace('\tr.GET("/meta/brands", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetBrandDashboard)\n', '')
write_file('internal/meta/ad_account/route.go', adacc_route)

# 5. core/brand route
brand_route = read_file('internal/core/brand/route.go')
brand_route = brand_route.replace('h.FindAll)\n', 'h.FindAll)\n\t\tg.GET("/dashboard", middleware.RequirePermission("core.brand.view"), h.GetBrandDashboard)\n')
write_file('internal/core/brand/route.go', brand_route)

print("Done")
