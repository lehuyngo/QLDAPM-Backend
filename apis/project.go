package apis

type Project struct {
	UUID                   string     `json:"uuid,omitempty"`
	FullName               string     `json:"fullname,omitempty"`
	ShortName              string     `json:"shortname,omitempty"`
	Code                   string     `json:"code,omitempty"`
	LastActiveTime         int64      `json:"last_active_time,omitempty"`
	LastMeetingCreatedTime int64      `json:"last_meeting_created_time,omitempty"`
	CreatedTime            int64      `json:"created_time,omitempty"`
	ProjectStatus          int32      `json:"project_status,omitempty"`
	Client                 *Client    `json:"client,omitempty"`
	Contacts               []*Contact `json:"contacts,omitempty"`
}

type CreateProjectRequest struct {
	FullName      string  `json:"fullname,omitempty"`
	ShortName     string  `json:"shortname,omitempty"`
	Code          string  `json:"code,omitempty"`
	ProjectStatus int32   `json:"project_status,omitempty"`
	Client        Client  `json:"client,omitempty"`
	Contact       Contact `json:"contact,omitempty"`
}

type CreateProjectStatusRequest struct {
	ProjectStatus int32 `json:"project_status,omitempty"`
}

func (r *CreateProjectRequest) IsClientValid() bool {
	return (r.Client.UUID != "") || (r.Client.FullName != "")
}

func (r *CreateProjectRequest) IsContactValid() bool {
	return (r.Contact.UUID != "") || (r.Contact.FullName != "")
}

type ListProjectResponse struct {
	Data []*Project `json:"data"`
}

type ProjectStatus struct {
	Value int32  `json:"value,omitempty"`
	Name  string `json:"name,omitempty"`
	Note  string `json:"note,omitempty"`
}

type ListProjectStatusResponse struct {
	Data []*ProjectStatus `json:"data"`
}
