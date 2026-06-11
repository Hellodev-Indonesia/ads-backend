package ads

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultAdFields = "id,campaign_id,adset_id,name,status,effective_status,creative,created_time,updated_time"
const DefaultCreativeFields = "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,url_tags"

type Service interface {
	// DB reads (used by handlers)
	GetAds(filter AdFilter) ([]dto.AdResponse, *response.PaginationMeta, error)

	// Direct Meta API call (creatives stay as direct calls)
	GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error)

	// Meta API sync (used by sync job)
	SyncAds(adAccountID string) (int, error)
	SyncAdsWithList(adAccountID string) (int, []MetaAd, error)
	SyncAdsByIDs(adAccountID string, ids []string) (int, []MetaAd, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS ---

func (s *serviceImpl) GetAds(filter AdFilter) ([]dto.AdResponse, *response.PaginationMeta, error) {
	adsList, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ads: %w", err)
	}

	var result []dto.AdResponse
	for _, a := range adsList {
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

// --- LOCAL DB READ (creatives are synced in background) ---

func (s *serviceImpl) GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error) {
	// The fields parameter is ignored since we return the exact JSON payload
	// that was synced to the local database in the background job.
	payload, err := s.repo.FindCreativeRawPayload(creativeID)
	if err != nil {
		return nil, fmt.Errorf("creative not found in database: %w", err)
	}

	var item dto.CreativeResponse
	if err := json.Unmarshal([]byte(payload), &item); err != nil {
		return nil, fmt.Errorf("failed to parse creative JSON from db: %w", err)
	}

	return &item, nil
}

// --- META API SYNC ---

func (s *serviceImpl) SyncAds(adAccountID string) (int, error) {
	count, _, err := s.SyncAdsWithList(adAccountID)
	return count, err
}

func (s *serviceImpl) SyncAdsWithList(adAccountID string) (int, []MetaAd, error) {
	params := url.Values{}
	params.Set("fields", DefaultAdFields)

	rawList, _, err := s.client.Get(adAccountID+"/ads", params, true)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to fetch ads from Meta: %w", err)
	}

	var models []MetaAd
	for _, raw := range rawList {
		var item dto.AdResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping ad unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, nil, fmt.Errorf("failed to upsert ads: %w", err)
	}

	return len(models), models, nil
}

func (s *serviceImpl) SyncAdsByIDs(adAccountID string, ids []string) (int, []MetaAd, error) {
	if len(ids) == 0 {
		return 0, nil, nil
	}

	batchSize := 50
	var allModels []MetaAd

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[i:end]

		params := url.Values{}
		params.Set("ids", strings.Join(batchIDs, ","))
		params.Set("fields", DefaultAdFields)

		rawList, _, err := s.client.Get("", params, false)
		if err != nil {
			log.Printf("Warning: failed to bulk fetch ads: %v", err)
			continue
		}

		if len(rawList) == 0 {
			continue
		}

		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(rawList[0], &rawMap); err != nil {
			log.Printf("Warning: failed to unmarshal bulk ads response: %v", err)
			continue
		}

		for _, rawPayloadBytes := range rawMap {
			var item dto.AdResponse
			if err := json.Unmarshal(rawPayloadBytes, &item); err != nil {
				continue
			}
			allModels = append(allModels, mapDTOToModel(item))
		}
	}

	if len(allModels) > 0 {
		if err := s.repo.UpsertBatch(allModels); err != nil {
			return 0, nil, fmt.Errorf("failed to upsert ads: %w", err)
		}
	}

	return len(allModels), allModels, nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaAd) dto.AdResponse {
	return dto.AdResponse{
		ID:              m.ID,
		Name:            m.Name,
		AdSetID:         m.AdSetID,
		CampaignID:      m.CampaignID,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		Creative:        dto.CreativeRef{ID: m.CreativeID},
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.AdResponse) MetaAd {
	return MetaAd{
		ID:              d.ID,
		CampaignID:      d.CampaignID,
		AdSetID:         d.AdSetID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		CreativeID:      d.Creative.ID,
		CreatedTime:     parseTime(d.CreatedTime),
		UpdatedTime:     parseTime(d.UpdatedTime),
	}
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
