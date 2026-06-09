package ad_account

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

type Service interface {
	GetAdAccounts(filter AdAccountFilter) ([]dto.AdAccountResponse, *response.Meta, error)
	GetUnassigned(filter AdAccountFilter) ([]dto.AdAccountResponse, *response.Meta, error)
	BulkAssignBrand(req dto.AssignBrandRequest) error
	SyncAdAccounts() (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

func (s *serviceImpl) GetAdAccounts(filter AdAccountFilter) ([]dto.AdAccountResponse, *response.Meta, error) {
	accounts, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.AdAccountResponse
	for _, acc := range accounts {
		result = append(result, dto.AdAccountResponse{
			ID:            acc.ID,
			Name:          acc.Name,
			AccountStatus: acc.AccountStatus,
			BrandID:       acc.BrandID,
			Currency:      acc.Currency,
			TimezoneName:  acc.TimezoneName,
			BusinessID:    acc.BusinessID,
			BusinessName:  acc.BusinessName,
			IsActive:      acc.IsActive,
		})
	}

	lastPage := int(total) / filter.Limit
	if int(total)%filter.Limit != 0 {
		lastPage++
	}

	meta := &response.Meta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    total,
		LastPage: lastPage,
	}

	return result, meta, nil
}

func (s *serviceImpl) GetUnassigned(filter AdAccountFilter) ([]dto.AdAccountResponse, *response.Meta, error) {
	accounts, total, err := s.repo.FindUnassigned(filter)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.AdAccountResponse
	for _, acc := range accounts {
		result = append(result, dto.AdAccountResponse{
			ID:            acc.ID,
			Name:          acc.Name,
			AccountStatus: acc.AccountStatus,
			BrandID:       acc.BrandID,
			Currency:      acc.Currency,
			TimezoneName:  acc.TimezoneName,
			BusinessID:    acc.BusinessID,
			BusinessName:  acc.BusinessName,
			IsActive:      acc.IsActive,
		})
	}

	lastPage := int(total) / filter.Limit
	if filter.Limit > 0 && int(total)%filter.Limit != 0 {
		lastPage++
	}

	meta := &response.Meta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    total,
		LastPage: lastPage,
	}

	return result, meta, nil
}

func (s *serviceImpl) BulkAssignBrand(req dto.AssignBrandRequest) error {
	// We rely on database Foreign Key constraint to validate if BrandID exists
	// This avoids circular dependency between meta and core domains.
	if len(req.AdAccountIDs) > 0 {
		return s.repo.UpdateBrandIDBatch(req.AdAccountIDs, req.BrandID)
	} else if req.BusinessID != nil {
		return s.repo.UpdateBrandIDByBusiness(*req.BusinessID, req.BrandID)
	}
	return nil
}

func (s *serviceImpl) SyncAdAccounts() (int, error) {
	params := url.Values{}
	params.Set("fields", "id,name,account_status,currency,timezone_name,business")

	rawList, _, err := s.client.Get("me/adaccounts", params, true)
	if err != nil {
		return 0, err
	}

	var models []MetaAdAccount
	for _, raw := range rawList {
		var item struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			AccountStatus int    `json:"account_status"`
			Currency      string `json:"currency"`
			TimezoneName  string `json:"timezone_name"`
			Business      struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"business"`
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			return 0, err
		}

		var currency *string
		if item.Currency != "" {
			currency = &item.Currency
		}

		var timezone *string
		if item.TimezoneName != "" {
			timezone = &item.TimezoneName
		}

		var businessID *string
		var businessName *string
		if item.Business.ID != "" {
			businessID = &item.Business.ID
			businessName = &item.Business.Name
		}

		models = append(models, MetaAdAccount{
			ID:            item.ID,
			Name:          item.Name,
			AccountStatus: item.AccountStatus,
			Currency:      currency,
			TimezoneName:  timezone,
			BusinessID:    businessID,
			BusinessName:  businessName,
			IsActive:      item.AccountStatus == 1,
		})
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, err
	}

	return len(models), nil
}
