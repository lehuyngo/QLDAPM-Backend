package apis

type ContactActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	Contact     *Contact     `json:"contact,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityContact struct {
	Data []*ContactActivity `json:"data,omitempty"`
}
