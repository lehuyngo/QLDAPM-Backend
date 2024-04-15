package apis

type ClientActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	Client      *Client      `json:"client,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityClient struct {
	Data []*ClientActivity `json:"data,omitempty"`
}
