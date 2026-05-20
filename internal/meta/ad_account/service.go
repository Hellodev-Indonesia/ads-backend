package ad_account

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	GetAdAccounts() ([]dto.AdAccountResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAdAccounts() ([]dto.AdAccountResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,account_status")

	rawList, err := s.client.Get("me/adaccounts", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.AdAccountResponse
	for _, raw := range rawList {
		var item dto.AdAccountResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
