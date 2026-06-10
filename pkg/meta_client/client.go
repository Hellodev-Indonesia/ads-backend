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
	Previous string `json:"previous"`
	Next     string `json:"next"`
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
			Timeout: 60 * time.Second,
		},
	}
}

// Get performs a GET request, attaches the access_token, handles paging (if autoPage is true), and parses errors.
func (c *Client) Get(path string, queryParams url.Values, autoPage bool) ([]json.RawMessage, *Paging, error) {
	if c.AccessToken == "" {
		return nil, nil, errors.New("Meta Access Token is missing. Please set META_ACCESS_TOKEN in your environment/.env")
	}

	fullURL := fmt.Sprintf("%s/%s/%s", c.BaseURL, c.Version, path)

	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse URL: %w", err)
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
	var lastPaging *Paging

	for nextURL != "" {
		var resp *http.Response
		var bodyBytes []byte
		var reqErr error

		maxRetries := 0
		for attempt := 0; attempt <= maxRetries; attempt++ {
			req, err := http.NewRequest("GET", nextURL, nil)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create request: %w", err)
			}

			resp, err = c.HTTPClient.Do(req)
			if err != nil {
				reqErr = fmt.Errorf("failed to execute request: %w", err)
				log.Printf("Meta API Request failed: GET %s (attempt %d) - Error: %v", path, attempt+1, err)
			} else {
				bodyBytes, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					reqErr = fmt.Errorf("failed to read response body: %w", err)
				} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					var wrap errorWrapper
					isRateLimit := false
					if jsonErr := json.Unmarshal(bodyBytes, &wrap); jsonErr == nil && wrap.Error != nil {
						// 17: User limit, 4: App limit, 80004: Ads API limit, 613: Custom limit, 32: Page limit, 1/2: Temporary errors
						if wrap.Error.Code == 17 || wrap.Error.Code == 4 || wrap.Error.Code == 80004 || wrap.Error.Code == 613 || wrap.Error.Code == 32 || wrap.Error.Code == 1 || wrap.Error.Code == 2 {
							isRateLimit = true
						}
					}

					if isRateLimit {
						reqErr = wrap.Error
						log.Printf("Meta API Rate limit hit: %v", reqErr)
					} else {
						// Not a rate limit error, don't retry
						if wrap.Error != nil {
							return nil, nil, wrap.Error
						}
						return nil, nil, fmt.Errorf("Meta API returned status %d: %s", resp.StatusCode, string(bodyBytes))
					}
				} else {
					// Success
					reqErr = nil
					break
				}
			}

			if attempt < maxRetries {
				sleepDuration := time.Duration(1<<attempt) * 10 * time.Second // 10s, 20s, 40s, 80s, 160s
				log.Printf("Retrying Meta API request in %v...", sleepDuration)
				time.Sleep(sleepDuration)
			}
		}

		if reqErr != nil {
			return nil, nil, reqErr
		}

		var base BaseResponse
		if err := json.Unmarshal(bodyBytes, &base); err != nil {
			// Response might be a single object, not a paginated list of data
			return []json.RawMessage{bodyBytes}, nil, nil
		}

		if base.Data == nil && !autoPage {
			return []json.RawMessage{bodyBytes}, nil, nil
		}

		allData = append(allData, base.Data...)
		lastPaging = base.Paging

		if autoPage && base.Paging != nil && base.Paging.Next != "" {
			nextURL = base.Paging.Next
		} else {
			nextURL = ""
		}
	}

	return allData, lastPaging, nil
}
