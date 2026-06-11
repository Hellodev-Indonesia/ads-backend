package centrifugo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	url    string
	apiKey string
	http   *http.Client
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		url:    url,
		apiKey: apiKey,
		http:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Publish(ctx context.Context, channel string, data any) error {
	payload := map[string]any{
		"method": "publish",
		"params": map[string]any{
			"channel": channel,
			"data":    data,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/api", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		log.Printf("[Centrifugo] Publish HTTP Error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[Centrifugo] Publish Bad Status %d: %s", resp.StatusCode, string(bodyBytes))
		return fmt.Errorf("centrifugo returned status %d", resp.StatusCode)
	}

	dataStr := ""
	if b, err := json.Marshal(data); err == nil {
		if len(b) > 200 {
			dataStr = string(b[:197]) + "..."
		} else {
			dataStr = string(b)
		}
	}

	log.Printf("[Centrifugo] Successfully published to channel %s: %s", channel, dataStr)
	return nil
}
