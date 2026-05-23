package ad_creative

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule"
	fraudlog "github.com/alex/ads_backend/internal/core/fraud_log"
	fraudlogDto "github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"github.com/alex/ads_backend/internal/meta/ad_account"
	adsDto "github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/internal/notification/alert"
	alertDto "github.com/alex/ads_backend/internal/notification/alert/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"gorm.io/gorm"
)

const creativeFields = "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,url_tags"

// Typed structs for navigating object_story_spec nested JSON.
// object_story_spec comes back as interface{} from CreativeResponse, so we
// round-trip through JSON to get typed access.
type objectStorySpec struct {
	LinkData     *linkData     `json:"link_data"`
	VideoData    *videoData    `json:"video_data"`
	TemplateData *templateData `json:"template_data"`
}

type linkData struct {
	Link string `json:"link"`
}

type videoData struct {
	CallToAction *callToAction `json:"call_to_action"`
}

type callToAction struct {
	Value *callToActionValue `json:"value"`
}

type callToActionValue struct {
	Link string `json:"link"`
}

type templateData struct {
	Link string `json:"link"`
}

type assetFeedSpec struct {
	LinkURLs []struct {
		WebsiteURL string `json:"website_url"`
	} `json:"link_urls"`
}

type Service interface {
	SyncCreatives(adAccountID string, adsList []AdRecord) (int, error)
}

// AdRecord carries the minimal ad data needed for creative lookup.
type AdRecord struct {
	ID         string
	CreativeID string
	AdSetID    string
	CampaignID string
}

type serviceImpl struct {
	client        *meta_client.Client
	repo          Repository
	adAccountRepo ad_account.Repository
	whitelistSvc  brand_whitelist_rule.Service
	fraudLogSvc   fraudlog.Service
	alertSvc      alert.Service
}

func NewService(
	client *meta_client.Client,
	repo Repository,
	adAccountRepo ad_account.Repository,
	whitelistSvc brand_whitelist_rule.Service,
	fraudLogSvc fraudlog.Service,
	alertSvc alert.Service,
) Service {
	return &serviceImpl{
		client:        client,
		repo:          repo,
		adAccountRepo: adAccountRepo,
		whitelistSvc:  whitelistSvc,
		fraudLogSvc:   fraudLogSvc,
		alertSvc:      alertSvc,
	}
}

func (s *serviceImpl) SyncCreatives(adAccountID string, adsList []AdRecord) (int, error) {
	account, err := s.adAccountRepo.FindByID(adAccountID)
	if err != nil {
		return 0, fmt.Errorf("ad account not found: %w", err)
	}

	// Deduplicate creative IDs; keep the first ad per creative for context.
	seen := make(map[string]AdRecord)
	for _, ad := range adsList {
		if ad.CreativeID == "" {
			continue
		}
		if _, exists := seen[ad.CreativeID]; !exists {
			seen[ad.CreativeID] = ad
		}
	}

	count := 0
	for creativeID, ad := range seen {
		if err := s.processCreative(creativeID, ad, adAccountID, account.BrandID); err != nil {
			log.Printf("Warning: failed to process creative %s: %v", creativeID, err)
			continue
		}
		count++
	}
	return count, nil
}

func (s *serviceImpl) processCreative(creativeID string, ad AdRecord, adAccountID string, brandID *uint64) error {
	params := url.Values{}
	params.Set("fields", creativeFields)

	rawList, _, err := s.client.Get(creativeID, params, false)
	if err != nil {
		return fmt.Errorf("Meta API error: %w", err)
	}
	if len(rawList) == 0 {
		return nil
	}

	// Reuse existing ads/dto.CreativeResponse for top-level unmarshal.
	var apiResp adsDto.CreativeResponse
	if err := json.Unmarshal(rawList[0], &apiResp); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	destinationURL := extractDestinationURL(&apiResp)
	normalizedDomain := ""
	urlHash := ""
	if destinationURL != "" {
		normalizedDomain = extractHostname(destinationURL)
		urlHash = hashURL(destinationURL)
	}

	rawPayload := string(rawList[0])
	now := time.Now()

	// Load previous state before upsert.
	prev, prevErr := s.repo.FindByCreativeID(creativeID)
	isFirstSync := prevErr == gorm.ErrRecordNotFound

	creative := buildCreativeModel(creativeID, &apiResp, destinationURL, normalizedDomain, urlHash, rawPayload, now)

	if err := s.repo.Upsert(creative); err != nil {
		return fmt.Errorf("upsert creative: %w", err)
	}

	changeType, changedFields := detectChange(prev, creative, isFirstSync)
	if changeType == "unchanged" {
		return nil
	}

	version := buildVersion(creative, ad, adAccountID, brandID, changeType, changedFields, now)
	if err := s.repo.CreateVersion(version); err != nil {
		log.Printf("Warning: failed to save creative version for %s: %v", creativeID, err)
	}

	if brandID == nil || destinationURL == "" {
		return nil
	}

	s.runPolicyCheck(creative, version, *brandID, isFirstSync, changeType)
	return nil
}

func (s *serviceImpl) runPolicyCheck(
	creative *AdCreative,
	version *AdCreativeVersion,
	brandID uint64,
	isFirstSync bool,
	changeType string,
) {
	targetURL := *creative.DestinationURL

	result, err := s.whitelistSvc.IsURLAllowed(brandID, targetURL, "destination_url")
	if err != nil {
		log.Printf("Warning: whitelist check failed for creative %s: %v", creative.CreativeID, err)
		return
	}

	eventType, severity := classifyEvent(isFirstSync, changeType, result.Allowed)
	if eventType == "" {
		return // no fraud event
	}

	// Resolve old URL from previous version for the log.
	var oldURL *string
	if prev, err := s.repo.FindLatestVersionByCreativeID(creative.CreativeID); err == nil && prev.DestinationURL != nil {
		oldURL = prev.DestinationURL
	}

	newVal := targetURL
	if eventType == "creative_url_removed" {
		newVal = ""
	}

	// Dedup: don't create a second open fraud log for the same (creative, event, url).
	exists, err := s.fraudLogSvc.ExistsOpenDuplicate(creative.CreativeID, eventType, newVal)
	if err != nil {
		log.Printf("Warning: dedup check failed: %v", err)
	}
	if exists {
		return
	}

	var matchedRuleID *uint64
	if result.Rule != nil {
		matchedRuleID = &result.Rule.ID
	}

	msg := buildMessage(eventType, creative.CreativeID, targetURL)
	input := fraudlogDto.CreateFraudLogInput{
		BrandID:       &brandID,
		AdAccountID:   version.AdAccountID,
		CampaignID:    version.CampaignID,
		AdsetID:       version.AdsetID,
		AdID:          version.AdID,
		CreativeID:    &creative.CreativeID,
		EventType:     eventType,
		Severity:      severity,
		OldValue:      oldURL,
		NewValue:      &newVal,
		MatchedRuleID: matchedRuleID,
		Message:       msg,
	}

	fraudLog, err := s.fraudLogSvc.Create(input)
	if err != nil {
		log.Printf("Warning: failed to create fraud log: %v", err)
		return
	}

	_, err = s.alertSvc.Create(alertDto.CreateAlertInput{
		FraudLogID: &fraudLog.ID,
		BrandID:    &brandID,
		Title:      buildAlertTitle(eventType),
		Message:    msg,
		Severity:   severity,
	})
	if err != nil {
		log.Printf("Warning: failed to create alert: %v", err)
	}
}

// --- URL extraction ---

// extractDestinationURL navigates the Meta creative response to find the primary destination URL.
// object_story_spec is interface{} in CreativeResponse, so we round-trip through JSON.
func extractDestinationURL(c *adsDto.CreativeResponse) string {
	if c.ObjectStorySpec != nil {
		b, _ := json.Marshal(c.ObjectStorySpec)
		var spec objectStorySpec
		if err := json.Unmarshal(b, &spec); err == nil {
			if spec.LinkData != nil && spec.LinkData.Link != "" {
				return spec.LinkData.Link
			}
			if spec.VideoData != nil &&
				spec.VideoData.CallToAction != nil &&
				spec.VideoData.CallToAction.Value != nil &&
				spec.VideoData.CallToAction.Value.Link != "" {
				return spec.VideoData.CallToAction.Value.Link
			}
			if spec.TemplateData != nil && spec.TemplateData.Link != "" {
				return spec.TemplateData.Link
			}
		}
	}
	if c.AssetFeedSpec != nil {
		b, _ := json.Marshal(c.AssetFeedSpec)
		var feed assetFeedSpec
		if err := json.Unmarshal(b, &feed); err == nil && len(feed.LinkURLs) > 0 {
			if u := feed.LinkURLs[0].WebsiteURL; u != "" {
				return u
			}
		}
	}
	return ""
}

func extractHostname(raw string) string {
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}

func hashURL(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)[:16]
}

// --- Model builders ---

func buildCreativeModel(
	creativeID string,
	api *adsDto.CreativeResponse,
	destURL, normalizedDomain, urlHash, rawPayload string,
	now time.Time,
) *AdCreative {
	c := &AdCreative{
		CreativeID: creativeID,
		SyncedAt:   &now,
		RawPayload: &rawPayload,
	}
	if api.Name != "" {
		c.Name = &api.Name
	}
	if api.Title != "" {
		c.Title = &api.Title
	}
	if api.Body != "" {
		c.Body = &api.Body
	}
	if api.ImageURL != "" {
		c.ImageURL = &api.ImageURL
	}
	if destURL != "" {
		c.DestinationURL = &destURL
		c.NormalizedDomain = &normalizedDomain
		c.URLHash = &urlHash
	}
	return c
}

func buildVersion(
	creative *AdCreative,
	ad AdRecord,
	adAccountID string,
	brandID *uint64,
	changeType string,
	changedFields []string,
	now time.Time,
) *AdCreativeVersion {
	v := &AdCreativeVersion{
		CreativeID:       creative.CreativeID,
		AdID:             &ad.ID,
		AdsetID:          &ad.AdSetID,
		CampaignID:       &ad.CampaignID,
		AdAccountID:      &adAccountID,
		BrandID:          brandID,
		Name:             creative.Name,
		Title:            creative.Title,
		Body:             creative.Body,
		ImageURL:         creative.ImageURL,
		DestinationURL:   creative.DestinationURL,
		NormalizedDomain: creative.NormalizedDomain,
		URLHash:          creative.URLHash,
		RawPayload:       creative.RawPayload,
		ChangeType:       &changeType,
		SyncedAt:         &now,
	}
	if len(changedFields) > 0 {
		b, _ := json.Marshal(changedFields)
		s := string(b)
		v.ChangedFields = &s
	}
	return v
}

// --- Change detection ---

func detectChange(prev *AdCreative, curr *AdCreative, isFirstSync bool) (string, []string) {
	if isFirstSync {
		return "initial_sync", nil
	}

	prevURL := strVal(prev.DestinationURL)
	currURL := strVal(curr.DestinationURL)

	if prevURL != currURL {
		changed := []string{"destination_url"}
		switch {
		case prevURL == "":
			return "url_added", changed
		case currURL == "":
			return "url_removed", changed
		default:
			return "url_changed", changed
		}
	}

	var changed []string
	if prev != nil {
		if strVal(prev.Title) != strVal(curr.Title) {
			changed = append(changed, "title")
		}
		if strVal(prev.Body) != strVal(curr.Body) {
			changed = append(changed, "body")
		}
		if strVal(prev.ImageURL) != strVal(curr.ImageURL) {
			changed = append(changed, "image_url")
		}
	}
	if len(changed) > 0 {
		return "content_changed", changed
	}
	return "unchanged", nil
}

// classifyEvent maps (isFirstSync, changeType, allowed) to (eventType, severity).
// Returns ("", "") when no fraud event should be emitted.
func classifyEvent(isFirstSync bool, changeType string, allowed bool) (string, string) {
	switch {
	case isFirstSync && !allowed:
		return "destination_url_not_whitelisted", "high"
	case isFirstSync && allowed:
		return "", "" // first sync, whitelisted — no alert
	case changeType == "url_changed" && allowed:
		return "destination_url_changed_whitelisted", "low"
	case changeType == "url_changed" && !allowed:
		return "destination_url_changed_not_whitelisted", "high"
	case changeType == "url_removed":
		return "creative_url_removed", "medium"
	case changeType == "url_added" && !allowed:
		return "destination_url_not_whitelisted", "high"
	case changeType == "url_added" && allowed:
		return "", ""
	case !allowed:
		// Unchanged URL but still not whitelisted.
		return "destination_url_not_whitelisted", "high"
	default:
		return "", ""
	}
}

func buildMessage(eventType, creativeID, targetURL string) string {
	switch eventType {
	case "destination_url_not_whitelisted":
		return fmt.Sprintf("Creative %s has a non-whitelisted destination URL: %s", creativeID, targetURL)
	case "destination_url_changed_whitelisted":
		return fmt.Sprintf("Creative %s destination URL changed to a whitelisted URL: %s", creativeID, targetURL)
	case "destination_url_changed_not_whitelisted":
		return fmt.Sprintf("Creative %s destination URL changed to a non-whitelisted URL: %s", creativeID, targetURL)
	case "creative_url_removed":
		return fmt.Sprintf("Creative %s destination URL was removed", creativeID)
	default:
		return fmt.Sprintf("Creative %s policy event: %s", creativeID, eventType)
	}
}

func buildAlertTitle(eventType string) string {
	switch eventType {
	case "destination_url_not_whitelisted":
		return "Non-whitelisted URL Detected"
	case "destination_url_changed_whitelisted":
		return "URL Changed (Whitelisted)"
	case "destination_url_changed_not_whitelisted":
		return "URL Changed to Non-whitelisted"
	case "creative_url_removed":
		return "Creative URL Removed"
	default:
		return "Creative Policy Alert"
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
