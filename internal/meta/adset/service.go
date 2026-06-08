package adset

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultFields = "id,campaign_id,name,status,effective_status,daily_budget,lifetime_budget,budget_remaining,bid_strategy,attribution_spec,start_time,end_time,created_time,updated_time"

type Service interface {
	// DB reads (used by handlers)
	GetAdSets(filter AdSetFilter) ([]dto.AdSetResponse, *response.PaginationMeta, error)

	// Meta API sync (used by sync job)
	SyncAdSets(adAccountID string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS ---

func (s *serviceImpl) GetAdSets(filter AdSetFilter) ([]dto.AdSetResponse, *response.PaginationMeta, error) {
	adsets, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch adsets: %w", err)
	}

	var result []dto.AdSetResponse
	for _, a := range adsets {
		result = append(result, mapModelToDTO(a))
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

// --- META API SYNC ---

func (s *serviceImpl) SyncAdSets(adAccountID string) (int, error) {
	params := url.Values{}
	params.Set("fields", DefaultFields)

	rawList, _, err := s.client.Get(adAccountID+"/adsets", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch adsets from Meta: %w", err)
	}

	var models []MetaAdSet
	for _, raw := range rawList {
		var item dto.AdSetResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping adset unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert adsets: %w", err)
	}

	return len(models), nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaAdSet) dto.AdSetResponse {
	return dto.AdSetResponse{
		ID:              m.ID,
		Name:            m.Name,
		CampaignID:      m.CampaignID,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		DailyBudget:     formatDecimal(m.DailyBudget),
		LifetimeBudget:  formatDecimal(m.LifetimeBudget),
		BudgetRemaining: formatDecimal(m.BudgetRemaining),
		BidStrategy:     m.BidStrategy,
		AttributionSpec: m.AttributionSpec,
		StartTime:       formatTime(m.StartTime),
		EndTime:         formatTime(m.EndTime),
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.AdSetResponse) MetaAdSet {
	var attrSpec json.RawMessage
	if d.AttributionSpec != nil {
		attrSpec, _ = json.Marshal(d.AttributionSpec)
	}

	return MetaAdSet{
		ID:              d.ID,
		CampaignID:      d.CampaignID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		DailyBudget:     parseDecimal(d.DailyBudget),
		LifetimeBudget:  parseDecimal(d.LifetimeBudget),
		BudgetRemaining: parseDecimal(d.BudgetRemaining),
		BidStrategy:     d.BidStrategy,
		AttributionSpec: attrSpec,
		StartTime:       parseTime(d.StartTime),
		EndTime:         parseTime(d.EndTime),
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
