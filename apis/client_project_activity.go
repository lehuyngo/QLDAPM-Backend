package apis

type ClientProjectActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Project     *Project     `json:"project,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityClientProject struct {
	Data []*ClientProjectActivity `json:"data,omitempty"`
}
