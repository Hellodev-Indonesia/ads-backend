import re

def read_file(path):
    with open(path, 'r') as f: return f.read()

def write_file(path, content):
    with open(path, 'w') as f: f.write(content)

repo_content = read_file('internal/meta/dashboard/repository.go')
svc_content = read_file('internal/meta/dashboard/service.go')
hndlr_content = read_file('internal/meta/dashboard/handler.go')

brand_scan = re.search(r'(type brandDashboardScan struct \{.*?\n\})', repo_content, re.MULTILINE | re.DOTALL).group(1)
brand_repo = re.search(r'(func \(r \*repository\) FindBrandDashboard.*?^})', repo_content, re.MULTILINE | re.DOTALL).group(1)
brand_repo = brand_repo.replace('DashboardFilter', 'AdAccountFilter')

adacc_repo_file = read_file('internal/meta/ad_account/repository.go')
adacc_repo_file = adacc_repo_file.replace('DisconnectBrand(id string) error\n}', 'DisconnectBrand(id string) error\n\tFindBrandDashboard(filter AdAccountFilter) ([]brandDashboardScan, int64, error)\n}')
adacc_repo_file += '\n' + brand_scan + '\n' + brand_repo + '\n'
write_file('internal/meta/ad_account/repository.go', adacc_repo_file)

brand_svc = re.search(r'(func \(s \*serviceImpl\) GetBrandDashboard.*?^})', svc_content, re.MULTILINE | re.DOTALL).group(1)
brand_svc = brand_svc.replace('DashboardFilter', 'AdAccountFilter')

adacc_svc_file = read_file('internal/meta/ad_account/service.go')
adacc_svc_file = adacc_svc_file.replace('GetBusinessOptions() ([]dto.BusinessOptionResponse, error)\n}', 'GetBusinessOptions() ([]dto.BusinessOptionResponse, error)\n\tGetBrandDashboard(filter AdAccountFilter) ([]dto.BrandDashboardResponse, *response.PaginationMeta, error)\n}')
adacc_svc_file += '\n' + brand_svc + '\n'
write_file('internal/meta/ad_account/service.go', adacc_svc_file)

brand_hndlr = re.search(r'(// GetBrandDashboard godoc.*?^})', hndlr_content, re.MULTILINE | re.DOTALL).group(1)
brand_hndlr = brand_hndlr.replace('DashboardFilter', 'AdAccountFilter')

adacc_hndlr_file = read_file('internal/meta/ad_account/handler.go')
adacc_hndlr_file += '\n' + brand_hndlr + '\n'
write_file('internal/meta/ad_account/handler.go', adacc_hndlr_file)

route_file = read_file('internal/meta/ad_account/route.go')
route_file = route_file.replace('h.GetAdAccounts)', 'h.GetAdAccounts)\n\tr.GET("/meta/brands", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetBrandDashboard)')
write_file('internal/meta/ad_account/route.go', route_file)
