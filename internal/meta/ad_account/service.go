package ad_account

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultFields = "id,name,account_status"

type Service interface {
	GetAdAccounts(fields string, limit string, after string, before string, autoPage bool) ([]dto.AdAccountResponse, *response.MetaPaging, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAdAccounts(fields string, limit string, after string, before string, autoPage bool) ([]dto.AdAccountResponse, *response.MetaPaging, error) {
	params := url.Values{}
	params.Set("fields", fields)
	if limit != "" {
		params.Set("limit", limit)
	}
	if after != "" {
		params.Set("after", after)
	}
	if before != "" {
		params.Set("before", before)
	}

	rawList, paging, err := s.client.Get("me/adaccounts", params, autoPage)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.AdAccountResponse
	for _, raw := range rawList {
		var item dto.AdAccountResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, nil, err
		}
		result = append(result, item)
	}

	return result, mapPaging(paging), nil
}

func mapPaging(p *meta_client.Paging) *response.MetaPaging {
	if p == nil {
		return nil
	}
	res := &response.MetaPaging{}
	res.Cursors.Before = p.Cursors.Before
	res.Cursors.After = p.Cursors.After
	res.HasPrevious = p.Previous != ""
	res.HasNext = p.Next != ""
	return res
}
