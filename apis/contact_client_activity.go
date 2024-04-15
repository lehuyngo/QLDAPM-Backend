package apis

type ContactClientActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Client      *Client      `json:"client,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityContactClient struct {
	Data []*ContactClientActivity `json:"data,omitempty"`
}
