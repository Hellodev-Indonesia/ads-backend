package activity

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/alex/ads_backend/internal/meta/activity/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	FindAll(filter dto.ActivityFilter) ([]dto.ActivityResponse, int64, error)
	FindAllByBrand(brandID uint64, filter dto.ActivityFilter) ([]dto.ActivityResponse, int64, error)
	SyncActivities(adAccountID string) (int, error)
	FindLatestByObjectIDs(adAccountID string, objectIDs []string) (*MetaActivity, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client, repo}
}

func (s *serviceImpl) FindAll(filter dto.ActivityFilter) ([]dto.ActivityResponse, int64, error) {
	rows, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.ActivityResponse
	for _, r := range rows {
		results = append(results, toResponse(r))
	}
	return results, total, nil
}

func (s *serviceImpl) FindAllByBrand(brandID uint64, filter dto.ActivityFilter) ([]dto.ActivityResponse, int64, error) {
	rows, total, err := s.repo.FindAllByBrand(brandID, filter)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.ActivityResponse
	for _, r := range rows {
		results = append(results, toResponse(r))
	}
	return results, total, nil
}

func (s *serviceImpl) FindLatestByObjectIDs(adAccountID string, objectIDs []string) (*MetaActivity, error) {
	return s.repo.FindLatestByObjectIDs(adAccountID, objectIDs)
}

func toResponse(act ActivityWithAdAccount) dto.ActivityResponse {
	// Account name formatting
	acct := act.AdAccountID
	if act.AdAccountName != nil {
		acct = *act.AdAccountName + "\n" + acct
	}
	if act.AdAccountBusinessName != nil {
		acct = *act.AdAccountBusinessName + " - " + acct
	}

	var extraData interface{}
	if len(act.ExtraData) > 0 && string(act.ExtraData) != "null" {
		json.Unmarshal(act.ExtraData, &extraData)
	}

	eventTimeStr := ""
	if act.EventTime != nil {
		eventTimeStr = act.EventTime.Format("2006-01-02 15:04:05")
	}

	var brand *dto.SimpleBrand
	if act.BrandID != nil && act.BrandName != nil {
		brand = &dto.SimpleBrand{
			ID:    *act.BrandID,
			Name:  *act.BrandName,
			Photo: act.BrandPhoto,
		}
	}

	return dto.ActivityResponse{
		ID:          act.ID,
		AdAccountID: act.AdAccountID,
		AdAccount:   acct,
		Brand:       brand,
		ActorID:     act.ActorID,
		ActorName:   act.ActorName,
		ObjectID:    act.ObjectID,
		ObjectName:  act.ObjectName,
		ObjectType:  act.ObjectType,
		EventType:   act.EventType,
		EventTime:   &eventTimeStr,
		ExtraData:   extraData,
		CreatedAt:   act.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   act.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// sync struct for meta response
type metaActivityRaw struct {
	Id         string          `json:"id"`
	ActorId    *string         `json:"actor_id"`
	ActorName  *string         `json:"actor_name"`
	ObjectId   *string         `json:"object_id"`
	ObjectName *string         `json:"object_name"`
	ObjectType *string         `json:"object_type"`
	EventType  *string         `json:"event_type"`
	EventTime  *string         `json:"event_time"`
	ExtraData  json.RawMessage `json:"extra_data"`
}

func (s *serviceImpl) SyncActivities(adAccountID string) (int, error) {
	params := url.Values{}
	params.Set("fields", "actor_id,actor_name,object_id,object_name,object_type,event_type,event_time,extra_data")

	rawList, _, err := s.client.Get(adAccountID+"/activities", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch activities from Meta: %w", err)
	}

	var models []MetaActivity
	for _, raw := range rawList {
		var item metaActivityRaw
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping activity unmarshal error: %v", err)
			continue
		}

		var eventTime *time.Time
		if item.EventTime != nil {
			if t, err := time.Parse(time.RFC3339, *item.EventTime); err == nil {
				eventTime = &t
			}
		}

		models = append(models, MetaActivity{
			AdAccountID: adAccountID,
			ActorID:     item.ActorId,
			ActorName:   item.ActorName,
			ObjectID:    item.ObjectId,
			ObjectName:  item.ObjectName,
			ObjectType:  item.ObjectType,
			EventType:   item.EventType,
			EventTime:   eventTime,
			ExtraData:   item.ExtraData,
		})
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert activities: %w", err)
	}

	return len(models), nil
}
