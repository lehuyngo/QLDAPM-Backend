package apis

type ClientContactActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Contact     *Contact     `json:"contact,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityClientContact struct {
	Data []*ClientContactActivity `json:"data,omitempty"`
}
