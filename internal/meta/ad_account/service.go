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

func (s *serviceImpl) SyncAdAccounts() (int, error) {
	params := url.Values{}
	params.Set("fields", "id,name,account_status")
	
	rawList, _, err := s.client.Get("me/adaccounts", params, true)
	if err != nil {
		return 0, err
	}

	var models []MetaAdAccount
	for _, raw := range rawList {
		var item dto.AdAccountResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return 0, err
		}
		
		models = append(models, MetaAdAccount{
			ID:            item.ID,
			Name:          item.Name,
			AccountStatus: item.AccountStatus,
		})
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, err
	}

	return len(models), nil
}
