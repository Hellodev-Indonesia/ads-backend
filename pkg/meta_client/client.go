package meta_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL     string
	Version     string
	AccessToken string
	HTTPClient  *http.Client
}

type Error struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode"`
	FBTraceID    string `json:"fbtrace_id"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Meta API Error: [%s] Code %d, Subcode %d: %s (trace_id: %s)", e.Type, e.Code, e.ErrorSubcode, e.Message, e.FBTraceID)
}

type errorWrapper struct {
	Error *Error `json:"error"`
}

type Paging struct {
	Cursors struct {
		Before string `json:"before"`
		After  string `json:"after"`
	} `json:"cursors"`
	Next string `json:"next"`
}

type BaseResponse struct {
	Data   []json.RawMessage `json:"data"`
	Paging *Paging           `json:"paging"`
}

func NewClient(baseURL, version, accessToken string) *Client {
	return &Client{
		BaseURL:     baseURL,
		Version:     version,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Get performs a GET request, attaches the access_token, handles paging (if autoPage is true), and parses errors.
func (c *Client) Get(path string, queryParams url.Values, autoPage bool) ([]json.RawMessage, error) {
	if c.AccessToken == "" {
		return nil, errors.New("Meta Access Token is missing. Please set META_ACCESS_TOKEN in your environment/.env")
	}

	fullURL := fmt.Sprintf("%s/%s/%s", c.BaseURL, c.Version, path)
	
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	q := u.Query()
	for k, v := range queryParams {
		for _, val := range v {
			q.Add(k, val)
		}
	}
	q.Set("access_token", c.AccessToken)
	u.RawQuery = q.Encode()

	var allData []json.RawMessage
	nextURL := u.String()

	for nextURL != "" {
		req, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			log.Printf("Meta API Request failed: GET %s - Error: %v", path, err)
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Meta API returned status %d for GET %s", resp.StatusCode, path)
			
			var wrap errorWrapper
			if err := json.Unmarshal(bodyBytes, &wrap); err == nil && wrap.Error != nil {
				return nil, wrap.Error
			}
			
			return nil, fmt.Errorf("Meta API returned status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var base BaseResponse
		if err := json.Unmarshal(bodyBytes, &base); err != nil {
			// Response might be a single object, not a paginated list of data
			return []json.RawMessage{bodyBytes}, nil
		}

		if base.Data == nil && !autoPage {
			return []json.RawMessage{bodyBytes}, nil
		}

		allData = append(allData, base.Data...)

		if autoPage && base.Paging != nil && base.Paging.Next != "" {
			nextURL = base.Paging.Next
		} else {
			nextURL = ""
		}
	}

	return allData, nil
}
