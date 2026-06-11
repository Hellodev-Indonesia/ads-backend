import re
import os

def read_file(path):
    with open(path, 'r') as f: return f.read()

def write_file(path, content):
    with open(path, 'w') as f: f.write(content)

repo_content = read_file('internal/meta/dashboard/repository.go')
svc_content = read_file('internal/meta/dashboard/service.go')
hndlr_content = read_file('internal/meta/dashboard/handler.go')

# 1. CAMPAIGN
campaign_scan = re.search(r'(type campaignDashboardScan struct \{.*?^\})', repo_content, re.MULTILINE | re.DOTALL).group(1)
campaign_repo = re.search(r'(func \(r \*repository\) FindCampaignDashboard.*?^})', repo_content, re.MULTILINE | re.DOTALL).group(1)

campaign_svc = re.search(r'(func \(s \*serviceImpl\) GetCampaignDashboard.*?^})', svc_content, re.MULTILINE | re.DOTALL).group(1)
campaign_map_scan = re.search(r'(func mapScanToDTO\(r campaignDashboardScan\) dto.CampaignDashboardRow \{.*?^\})', svc_content, re.MULTILINE | re.DOTALL).group(1)

campaign_hndlr = re.search(r'(// GetCampaignDashboard godoc.*?^})', hndlr_content, re.MULTILINE | re.DOTALL).group(1)

# Write to campaign
camp_repo_file = read_file('internal/meta/campaign/repository.go')
camp_repo_file += '\n' + campaign_scan + '\n' + campaign_repo.replace('DashboardFilter', 'CampaignFilter')
write_file('internal/meta/campaign/repository.go', camp_repo_file)

