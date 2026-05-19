package meta_client

import (
	"fmt"
	"net/http"
)

type MetaClient struct {
	AccessToken string
	AppID       string
	AppSecret   string
	BaseURL     string
}

func NewMetaClient(token, appID, appSecret string) *MetaClient {
	return &MetaClient{
		AccessToken: token,
		AppID:       appID,
		AppSecret:   appSecret,
		BaseURL:     "https://graph.facebook.com/v19.0",
	}
}

func (m *MetaClient) GetCampaigns(adAccountID string) (interface{}, error) {
	url := fmt.Sprintf("%s/%s/campaigns?access_token=%s", m.BaseURL, adAccountID, m.AccessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response...
	return nil, nil
}
