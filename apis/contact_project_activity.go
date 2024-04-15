package apis

type ContactProjectActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Project     *Project     `json:"project,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type ListActivityContactProject struct {
	Data []*ContactProjectActivity `json:"data,omitempty"`
}
