package apis

type ContactMailActivity struct {
	Type        ActivityType `json:"type,omitempty"`
	Creator     *User        `json:"creator,omitempty"`
	CreatedTime int64        `json:"created_time,omitempty"`
	Mail        *Mail        `json:"mail,omitempty"`
	Contact     *Contact     `json:"contact,omitempty"`
}

type ListActivityContactMail struct {
	Data []*ContactMailActivity `json:"data,omitempty"`
}
