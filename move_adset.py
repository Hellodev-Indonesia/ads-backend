import re

def read_file(path):
    with open(path, 'r') as f: return f.read()

def write_file(path, content):
    with open(path, 'w') as f: f.write(content)

repo_content = read_file('internal/meta/dashboard/repository.go')
svc_content = read_file('internal/meta/dashboard/service.go')
hndlr_content = read_file('internal/meta/dashboard/handler.go')

adset_scan = re.search(r'(type adSetDashboardScan struct \{.*?^\})', repo_content, re.MULTILINE | re.DOTALL).group(1)
adset_repo = re.search(r'(func \(r \*repository\) FindAdSetDashboard.*?^})', repo_content, re.MULTILINE | re.DOTALL).group(1)
adset_repo = adset_repo.replace('DashboardFilter', 'AdSetFilter')

helpers_go = """
func javaStrToInt(s string) (int, error) {
	var v int
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v, nil
}

func intToStr(v int) string {
	return fmt.Sprintf("%d", v)
}
"""

camp_repo_file = read_file('internal/meta/adset/repository.go')
camp_repo_file = camp_repo_file.replace('"gorm.io/gorm"\n\t"gorm.io/gorm/clause"', '"encoding/json"\n\t"fmt"\n\t"time"\n\n\t"gorm.io/gorm"\n\t"gorm.io/gorm/clause"')
camp_repo_file = camp_repo_file.replace('Search     string', 'Search     string\n\tDateStart  string\n\tDateStop   string')
camp_repo_file = camp_repo_file.replace('FindByCampaignID(campaignID string) ([]MetaAdSet, error)\n}', 'FindByCampaignID(campaignID string) ([]MetaAdSet, error)\n\tFindAdSetDashboard(filter AdSetFilter) ([]adSetDashboardScan, int64, error)\n}')
camp_repo_file += '\n' + adset_scan + '\n' + adset_repo + '\n' + helpers_go
write_file('internal/meta/adset/repository.go', camp_repo_file)

adset_svc = re.search(r'(func \(s \*serviceImpl\) GetAdSetDashboard.*?^})', svc_content, re.MULTILINE | re.DOTALL).group(1)
adset_svc = adset_svc.replace('DashboardFilter', 'AdSetFilter')
adset_map_scan = re.search(r'(func mapAdSetScanToDTO\(r adSetDashboardScan\) dto.AdSetDashboardRow \{.*?^\})', svc_content, re.MULTILINE | re.DOTALL).group(1)

svc_helpers = """
type metaAction struct {
	ActionType string `json:"action_type"`
	Value      string `json:"value"`
}

func resolveBudget(daily, lifetime float64) string {
	if daily > 0 {
		return formatFloat(daily)
	}
	return formatFloat(lifetime)
}

func parseActions(raw json.RawMessage) []metaAction {
	if raw == nil {
		return nil
	}
	var actions []metaAction
	_ = json.Unmarshal(raw, &actions)
	return actions
}

func findAction(actions []metaAction, actionType string) string {
	for _, a := range actions {
		if a.ActionType == actionType {
			return a.Value
		}
	}
	return "0"
}

func formatNullFloat(v *float64) string {
	if v == nil {
		return "0"
	}
	return formatFloat(*v)
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

func formatNullInt(v *int64) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatInt(*v, 10)
}
"""

camp_svc_file = read_file('internal/meta/adset/service.go')
camp_svc_file = camp_svc_file.replace('"context"\n\t"fmt"\n\t"log"', '"encoding/json"\n\t"fmt"\n\t"log"\n\t"math"\n\t"strconv"')
camp_svc_file = camp_svc_file.replace('GetAdSetByID(id string) (*dto.AdSetResponse, error)\n', 'GetAdSetByID(id string) (*dto.AdSetResponse, error)\n\tGetAdSetDashboard(filter AdSetFilter) ([]dto.AdSetDashboardRow, *response.PaginationMeta, error)\n')
camp_svc_file += '\n' + adset_svc + '\n' + adset_map_scan + '\n' + svc_helpers
write_file('internal/meta/adset/service.go', camp_svc_file)

adset_hndlr = re.search(r'(// GetAdSetDashboard godoc.*?^})', hndlr_content, re.MULTILINE | re.DOTALL).group(1)
adset_hndlr = adset_hndlr.replace('DashboardFilter', 'AdSetFilter')

camp_hndlr_file = read_file('internal/meta/adset/handler.go')
camp_hndlr_file += '\n' + adset_hndlr + '\n'
write_file('internal/meta/adset/handler.go', camp_hndlr_file)

route_file = read_file('internal/meta/adset/route.go')
route_file = route_file.replace('h.GetAdSetsByCampaign)', 'h.GetAdSetsByCampaign)\n\tr.GET("/meta/adsets/dashboard", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetAdSetDashboard)')
write_file('internal/meta/adset/route.go', route_file)

