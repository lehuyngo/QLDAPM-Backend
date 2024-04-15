package apis

type ClientNoteActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
	Note        *ClientNote  `json:"note,omitempty"`
}

type ListActivityClientNote struct {
	Data []*ClientNoteActivity `json:"data,omitempty"`
}
