package apis

type ClientTagActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
	Tag         *ClientTag   `json:"tag,omitempty"`
}

type ListActivityClientTag struct {
	Data []*ClientTagActivity `json:"data,omitempty"`
}
