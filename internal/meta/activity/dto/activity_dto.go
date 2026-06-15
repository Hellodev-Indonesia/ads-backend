package dto

type ActivityResponse struct {
	ID          uint64      `json:"id"`
	AdAccountID string      `json:"ad_account_id"`
	AdAccount   string      `json:"ad_account"`
	Brand       *SimpleBrand `json:"brand,omitempty"`
	ActorID     *string     `json:"actor_id"`
	ActorName   *string     `json:"actor_name"`
	ObjectID    *string     `json:"object_id"`
	ObjectName  *string     `json:"object_name"`
	ObjectType  *string     `json:"object_type"`
	EventType   *string     `json:"event_type"`
	EventTime   *string     `json:"event_time"`
	ExtraData   interface{} `json:"extra_data"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type ActivityFilter struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type SimpleBrand struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Photo *string `json:"photo"`
}
