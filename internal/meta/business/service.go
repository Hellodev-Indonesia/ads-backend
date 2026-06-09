package business

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/alex/ads_backend/internal/meta/business/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

type Service interface {
	SyncBusinesses() (int, error)
	GetBusinesses(filter BusinessFilter) ([]dto.BusinessResponse, *response.Meta, error)
}

type service struct {
	repo       Repository
	metaClient *meta_client.Client
}

func NewService(repo Repository, metaClient *meta_client.Client) Service {
	return &service{
		repo:       repo,
		metaClient: metaClient,
	}
}

type MetaAPIResponse struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ProfilePictureURI string `json:"profile_picture_uri"`
	TimezoneID        int    `json:"timezone_id"`
	CreatedTime       string `json:"created_time"`
}

func (s *service) SyncBusinesses() (int, error) {
	businessIDs, err := s.repo.GetUniqueBusinessIDsFromAdAccounts()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch unique business IDs: %w", err)
	}

	if len(businessIDs) == 0 {
		return 0, nil
	}

	queryParams := url.Values{}
	queryParams.Add("fields", "id,name,profile_picture_uri,timezone_id,created_time")

	var businesses []MetaBusiness
	for _, businessID := range businessIDs {
		rawMessages, _, err := s.metaClient.Get("/"+businessID, queryParams, false)
		if err != nil {
			fmt.Printf("Warning: failed to fetch business %s from Meta API: %v\n", businessID, err)
			continue
		}

		if len(rawMessages) == 0 {
			continue
		}

		var apiResp MetaAPIResponse
		// When fetching a single node, metaClient.Get might return it wrapped in an array or as a single object depending on how it's handled. 
		// Actually, metaClient.Get returns []json.RawMessage. If it's a single object without "data" wrapper, isEdge=false handles it and returns a slice with 1 element.
		if err := json.Unmarshal(rawMessages[0], &apiResp); err != nil {
			fmt.Printf("Warning: failed to decode business %s data: %v\n", businessID, err)
			continue
		}

		b := MetaBusiness{
			ID:   apiResp.ID,
			Name: apiResp.Name,
		}

		if apiResp.ProfilePictureURI != "" {
			uri := apiResp.ProfilePictureURI
			b.ProfilePictureURI = &uri
		}

		if apiResp.TimezoneID != 0 {
			tz := apiResp.TimezoneID
			b.TimezoneID = &tz
		}

		if apiResp.CreatedTime != "" {
			if parsedTime, err := time.Parse(time.RFC3339, apiResp.CreatedTime); err == nil {
				b.CreatedTime = &parsedTime
			}
		}

		businesses = append(businesses, b)
	}

	if len(businesses) == 0 {
		return 0, nil
	}

	if err := s.repo.UpsertBatch(businesses); err != nil {
		return 0, fmt.Errorf("failed to save businesses to database: %w", err)
	}

	return len(businesses), nil
}

func (s *service) GetBusinesses(filter BusinessFilter) ([]dto.BusinessResponse, *response.Meta, error) {
	businesses, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch businesses: %w", err)
	}

	var dtos []dto.BusinessResponse
	for _, b := range businesses {
		dtos = append(dtos, dto.BusinessResponse{
			ID:                b.ID,
			Name:              b.Name,
			ProfilePictureURI: b.ProfilePictureURI,
			TimezoneID:        b.TimezoneID,
			CreatedTime:       b.CreatedTime,
			SyncedAt:          b.SyncedAt,
		})
	}

	lastPage := int(math.Ceil(float64(total) / float64(filter.Limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	meta := &response.Meta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    total,
		LastPage: lastPage,
	}

	return dtos, meta, nil
}
