package meta

import (
	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Client interface {
	GetClient() *meta_client.Client
}

type clientImpl struct {
	rawClient *meta_client.Client
}

func NewClient() Client {
	return &clientImpl{
		rawClient: meta_client.NewClient(
			config.MetaGraphBaseURL,
			config.MetaGraphVersion,
			config.MetaAccessToken,
		),
	}
}

func (c *clientImpl) GetClient() *meta_client.Client {
	return c.rawClient
}
