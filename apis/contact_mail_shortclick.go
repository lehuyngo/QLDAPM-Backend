package apis

type ContactMailShortClick struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	Contact     *Contact     `json:"contact,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
}

type CreateContactMailShortClickRequest struct {
	Content 	string `json:"content,omitempty"`
}

type ListContactMailShortClick struct {
	Data []*ContactMailShortClick `json:"data,omitempty"`
}
