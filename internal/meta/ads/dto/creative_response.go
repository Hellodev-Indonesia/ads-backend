package dto

type CreativeResponse struct {
	ID              string      `json:"id" example:"2387654321098"`
	Name            string      `json:"name" example:"Creative Promo 1"`
	Title           string      `json:"title" example:"Special Offer"`
	Body            string      `json:"body" example:"Get 20% off today!"`
	ImageURL        string      `json:"image_url" example:"https://scontent.xx.fbcdn.net/..."`
	ThumbnailURL    string      `json:"thumbnail_url" example:"https://scontent.xx.fbcdn.net/..."`
	ObjectStorySpec interface{} `json:"object_story_spec,omitempty"`
	AssetFeedSpec   interface{} `json:"asset_feed_spec,omitempty"`
	URLTags         string      `json:"url_tags,omitempty" example:"utm_source=facebook&utm_medium=cpc"`
}
