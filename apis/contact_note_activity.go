package apis

type ContactNoteActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
	Note        *ContactNote `json:"note,omitempty"`
}

type ListActivityContactNote struct {
	Data []*ContactNoteActivity `json:"data,omitempty"`
}
