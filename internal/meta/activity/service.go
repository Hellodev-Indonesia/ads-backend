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
	FindAllByBrand(brandID uint64, filter dto.ActivityFilter) ([]dto.ActivityResponse, int64, error)
	SyncActivities(adAccountID string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client, repo}
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

func toResponse(act ActivityWithAdAccount) dto.ActivityResponse {
	// Account name formatting
	acct := act.AdAccountID
	if act.AdAccountName != nil {
		acct = *act.AdAccountName + "\n" + acct
	}
	if act.AdAccountBusinessName != nil {
		acct = *act.AdAccountBusinessName + " - " + acct
	}

	// Activity details from extra_data
	details := ""
	if len(act.ExtraData) > 0 && string(act.ExtraData) != "null" {
		var extraMap map[string]interface{}
		if err := json.Unmarshal(act.ExtraData, &extraMap); err == nil {
			if oldVal, ok := extraMap["old_value"]; ok {
				if newVal, ok := extraMap["new_value"]; ok {
					details = fmt.Sprintf("From %v to %v", oldVal, newVal)
				}
			}
		}
	}

	eventTimeStr := ""
	if act.EventTime != nil {
		eventTimeStr = act.EventTime.Format("02 Jan at 15.04")
	}

	eventType := ""
	if act.EventType != nil {
		eventType = formatEventType(*act.EventType)
	}

	objName := ""
	if act.ObjectName != nil {
		objName = *act.ObjectName
	} else if act.ObjectID != nil {
		objName = *act.ObjectID
	}

	actorName := ""
	if act.ActorName != nil {
		actorName = *act.ActorName
	}

	return dto.ActivityResponse{
		AdAccount:       acct,
		Activity:        eventType,
		ActivityDetails: details,
		ItemChanged:     objName,
		ChangeBy:        actorName,
		DateAndTime:     eventTimeStr,
	}
}

func formatEventType(e string) string {
	// Simple formatting, e.g. "update_status" -> "Update Status"
	// For production, a more comprehensive mapping is needed.
	switch e {
	case "update_status":
		return "Update Status"
	case "update_budget":
		return "Update Budget"
	case "create_ad_set":
		return "Create Ad Set"
	case "create_campaign":
		return "Create Campaign"
	case "create_ad":
		return "Create Ad"
	case "update_campaign_name":
		return "Update Campaign Name"
	default:
		return e
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
			ID:          item.Id,
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
