package apis

type ContactTagActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
	Tag         *ContactTag  `json:"tag,omitempty"`
}

type ListActivityContactTag struct {
	Data []*ContactTagActivity `json:"data,omitempty"`
}
