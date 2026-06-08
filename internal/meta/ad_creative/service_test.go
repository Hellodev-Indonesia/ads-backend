package ad_creative

import (
	"testing"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/stretchr/testify/assert"
)

func TestExtractHostname(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://www.example.com/path", "www.example.com"},
		{"http://example.com", "example.com"},
		{"example.com", "example.com"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractHostname(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractDestinationURL(t *testing.T) {
	tests := []struct {
		name     string
		input    *dto.CreativeResponse
		expected string
	}{
		{
			name: "Link Data",
			input: &dto.CreativeResponse{
				ObjectStorySpec: map[string]interface{}{
					"link_data": map[string]interface{}{
						"link": "https://example.com/link",
					},
				},
			},
			expected: "https://example.com/link",
		},
		{
			name: "Video Data",
			input: &dto.CreativeResponse{
				ObjectStorySpec: map[string]interface{}{
					"video_data": map[string]interface{}{
						"call_to_action": map[string]interface{}{
							"value": map[string]interface{}{
								"link": "https://example.com/video",
							},
						},
					},
				},
			},
			expected: "https://example.com/video",
		},
		{
			name: "Template Data",
			input: &dto.CreativeResponse{
				ObjectStorySpec: map[string]interface{}{
					"template_data": map[string]interface{}{
						"link": "https://example.com/template",
					},
				},
			},
			expected: "https://example.com/template",
		},
		{
			name: "Asset Feed Spec",
			input: &dto.CreativeResponse{
				AssetFeedSpec: map[string]interface{}{
					"link_urls": []map[string]interface{}{
						{"website_url": "https://example.com/feed"},
					},
				},
			},
			expected: "https://example.com/feed",
		},
		{
			name:     "Empty",
			input:    &dto.CreativeResponse{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDestinationURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClassifyEvent(t *testing.T) {
	tests := []struct {
		name        string
		isFirstSync bool
		changeType  string
		allowed     bool
		expEvent    string
		expSev      string
	}{
		{"First sync, not allowed", true, "initial_sync", false, "destination_url_not_whitelisted", "high"},
		{"First sync, allowed", true, "initial_sync", true, "", ""},
		{"URL changed, allowed", false, "url_changed", true, "destination_url_changed_whitelisted", "low"},
		{"URL changed, not allowed", false, "url_changed", false, "destination_url_changed_not_whitelisted", "high"},
		{"URL removed", false, "url_removed", true, "creative_url_removed", "medium"},
		{"URL added, not allowed", false, "url_added", false, "destination_url_not_whitelisted", "high"},
		{"URL added, allowed", false, "url_added", true, "", ""},
		{"Unchanged, not allowed", false, "unchanged", false, "destination_url_not_whitelisted", "high"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, sev := classifyEvent(tt.isFirstSync, tt.changeType, tt.allowed)
			assert.Equal(t, tt.expEvent, ev)
			assert.Equal(t, tt.expSev, sev)
		})
	}
}
