package apis

type Meeting struct {
	UUID           string         `json:"uuid,omitempty"`
	StartTime      int64          `json:"start_time,omitempty"`
	Location       string         `json:"location,omitempty"`
	Link           string         `json:"link,omitempty"`
	Attendees      []*Attendee    `json:"attendees,omitempty"`
	Creator        *User          `json:"creator,omitempty"`
	CreatedTime    int64          `json:"created_time,omitempty"`
	LastActiveTime int64          `json:"last_active_time,omitempty"`
	NoteCreators   []*NoteCreator `json:"note_creators,omitempty"`
}

type NoteCreator struct {
	Creator  *User  `json:"creator,omitempty"`
	NoteUUID string `json:"note_uuid,omitempty"`
}

type Attendee struct {
	UUID    string   `json:"uuid,omitempty"`
	Contact *Contact `json:"contact,omitempty"`
	User    *User    `json:"user,omitempty"`
}

type CreateMeetingRequest struct {
	StartTime    int64    `json:"start_time,omitempty"`
	Location     string   `json:"location,omitempty"`
	Link         string   `json:"link,omitempty"`
	UserUUIDs    []string `json:"user_uuids,omitempty"`
	ContactUUIDs []string `json:"contact_uuids,omitempty"`
}

type ListMeetingResponse struct {
	Data []*Meeting `json:"data"`
}
