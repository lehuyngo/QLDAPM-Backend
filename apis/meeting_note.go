package apis

type MeetingNote struct {
	UUID         string         `json:"uuid,omitempty"`
	Note         string         `json:"note,omitempty"`
	StartTime    int64          `json:"start_time,omitempty"`
	Location     string         `json:"location,omitempty"`
	Link         string         `json:"link,omitempty"`
	Contributors []*Contributor `json:"contributors,omitempty"`
	Editors      []*User        `json:"editors,omitempty"`
	Creator      *User          `json:"creator,omitempty"`
	CreatedTime  int64          `json:"created_time,omitempty"`
}

type Contributor struct {
	UUID    string   `json:"uuid,omitempty"`
	Contact *Contact `json:"contact,omitempty"`
	User    *User    `json:"user,omitempty"`
}

type CreateContributorRequest struct {
	ContactUUID string `json:"contact_uuid,omitempty"`
	UserUUID    string `json:"user_uuid,omitempty"`
}

type CreateContributorBatchRequest struct {
	ContactUUIDs []string `json:"contact_uuids,omitempty"`
	UserUUIDs    []string `json:"user_uuids,omitempty"`
}

type DeleteContributorBatchRequest struct {
	UUIDs []string `json:"uuids"`
}

type CreateMeetingNoteRequest struct {
	StartTime    int64    `json:"start_time,omitempty"`
	Location     string   `json:"location,omitempty"`
	Link         string   `json:"link,omitempty"`
	UserUUIDs    []string `json:"user_uuids,omitempty"`
	ContactUUIDs []string `json:"contact_uuids,omitempty"`
	Note         string   `json:"note,omitempty"`
}

type ListMeetingNoteResponse struct {
	Data []*MeetingNote `json:"data"`
}
