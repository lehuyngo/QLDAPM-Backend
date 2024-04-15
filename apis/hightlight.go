package apis

type Highlight struct {
	UUID            string `json:"uuid,omitempty"`
	Title           string `json:"title,omitempty"`
	MeetingNoteUUID string `json:"meeting_note_uuid,omitempty"`
	CreatedAt       int64  `json:"created_at,omitempty"`
	CreatorUUID     string `json:"creator_uuid,omitempty"`
}

type CreateMeetingHighlightRequest struct {
	Title []string `json:"titles"`
}

type CreateMeetingHighlightResponse struct {
	UUIDs []string `json:"uuids"`
}

type ListHighlightResponse struct {
	Data []*Highlight `json:"data"`
}

type DeleteMeetingHighlightRequest struct {
	UUIDs []string `json:"uuids"`
}
