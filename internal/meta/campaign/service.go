package campaign

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultFields = "id,name,status,effective_status,objective,daily_budget,lifetime_budget,budget_remaining,bid_strategy,start_time,stop_time,created_time,updated_time"

type Service interface {
	// DB reads (used by handlers)
	GetCampaigns(filter CampaignFilter) ([]dto.CampaignResponse, *response.PaginationMeta, error)
	GetCampaignByID(id string) (*dto.CampaignResponse, error)

	// Meta API sync (used by sync job)
	SyncCampaigns(adAccountID string) (int, error)
	SyncCampaignsByIDs(adAccountID string, ids []string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS (for API handlers) ---

func (s *serviceImpl) GetCampaigns(filter CampaignFilter) ([]dto.CampaignResponse, *response.PaginationMeta, error) {
	campaigns, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	var result []dto.CampaignResponse
	for _, c := range campaigns {
		result = append(result, mapModelToDTO(c))
	}

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	lastPage := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		lastPage++
	}

	meta := &response.PaginationMeta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    int(total),
		LastPage: lastPage,
	}

	return result, meta, nil
}

func (s *serviceImpl) GetCampaignByID(id string) (*dto.CampaignResponse, error) {
	campaign, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}
	result := mapModelToDTO(*campaign)
	return &result, nil
}

// --- META API SYNC METHOD (for background job) ---

func (s *serviceImpl) SyncCampaigns(adAccountID string) (int, error) {
	params := url.Values{}
	params.Set("fields", DefaultFields)

	rawList, _, err := s.client.Get(adAccountID+"/campaigns", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch campaigns from Meta: %w", err)
	}

	var models []MetaCampaign
	for _, raw := range rawList {
		var item dto.CampaignResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping campaign unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item, adAccountID))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert campaigns: %w", err)
	}

	return len(models), nil
}

func (s *serviceImpl) SyncCampaignsByIDs(adAccountID string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	batchSize := 50
	var allModels []MetaCampaign

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[i:end]

		params := url.Values{}
		params.Set("ids", strings.Join(batchIDs, ","))
		params.Set("fields", DefaultFields)

		rawList, _, err := s.client.Get("", params, false)
		if err != nil {
			log.Printf("Warning: failed to bulk fetch campaigns: %v", err)
			continue
		}

		if len(rawList) == 0 {
			continue
		}

		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(rawList[0], &rawMap); err != nil {
			log.Printf("Warning: failed to unmarshal bulk campaigns response: %v", err)
			continue
		}

		for _, rawPayloadBytes := range rawMap {
			var item dto.CampaignResponse
			if err := json.Unmarshal(rawPayloadBytes, &item); err != nil {
				continue
			}
			allModels = append(allModels, mapDTOToModel(item, adAccountID))
		}
	}

	if len(allModels) > 0 {
		if err := s.repo.UpsertBatch(allModels); err != nil {
			return 0, fmt.Errorf("failed to upsert campaigns: %w", err)
		}
	}

	return len(allModels), nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaCampaign) dto.CampaignResponse {
	return dto.CampaignResponse{
		ID:              m.ID,
		AccountID:       m.AccountID,
		Name:            m.Name,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		Objective:       m.Objective,
		BuyingType:      m.BuyingType,
		DailyBudget:     formatDecimal(m.DailyBudget),
		LifetimeBudget:  formatDecimal(m.LifetimeBudget),
		BudgetRemaining: formatDecimal(m.BudgetRemaining),
		SpendCap:        formatDecimal(m.SpendCap),
		BidStrategy:     m.BidStrategy,
		StartTime:       formatTime(m.StartTime),
		StopTime:        formatTime(m.StopTime),
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.CampaignResponse, accountID string) MetaCampaign {
	acctID := d.AccountID
	if acctID == "" {
		acctID = accountID
	}
	return MetaCampaign{
		ID:              d.ID,
		AccountID:       acctID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		Objective:       d.Objective,
		BuyingType:      d.BuyingType,
		DailyBudget:     parseDecimal(d.DailyBudget),
		LifetimeBudget:  parseDecimal(d.LifetimeBudget),
		BudgetRemaining: parseDecimal(d.BudgetRemaining),
		SpendCap:        parseDecimal(d.SpendCap),
		BidStrategy:     d.BidStrategy,
		StartTime:       parseTime(d.StartTime),
		StopTime:        parseTime(d.StopTime),
		CreatedTime:     parseTime(d.CreatedTime),
		UpdatedTime:     parseTime(d.UpdatedTime),
	}
}

func parseDecimal(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func formatDecimal(v float64) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	// Try multiple formats that Meta uses
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05+0700",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return &t
		}
	}
	return nil
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
