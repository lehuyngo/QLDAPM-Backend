package apis

type Client struct {
	UUID   			string `json:"uuid,omitempty"`
	FullName 		string `json:"fullname,omitempty"`
	ShortName 		string `json:"shortname,omitempty"`
	Code 			string `json:"code,omitempty"`
	Fax 			string `json:"fax,omitempty"`
	Website 		string `json:"website,omitempty"`
	Phone			string `json:"phone,omitempty"`
	Email			string `json:"email,omitempty"`
	CompanySize		string `json:"company_size,omitempty"`
	Address			string `json:"address,omitempty"`
	LastActiveTime	int64 `json:"last_active_time,omitempty"`
	CreatedTime		int64 `json:"created_time,omitempty"`
	Logo 			*File `json:"logo,omitempty"`
	Tags			[]*ClientTag `json:"tags,omitempty"`
	Notes			[]*ClientNote `json:"notes,omitempty"`
	Contacts		[]*Contact `json:"contacts,omitempty"`
	Projects		[]*Project `json:"projects,omitempty"`
}

func (c *Client) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

type CreateClientV1Request struct {
	Data Client `json:"data"`
	Force bool `json:"force,omitempty"`
}

type ListClientResponse struct {
	Data []*Client `json:"data"`
}

type CreateClientRequest struct {
	FullName 	string `form:"fullname,omitempty"`
	ShortName 	string `form:"shortname,omitempty"`
	Code 		string `form:"code,omitempty"`
	Fax 		string `form:"fax,omitempty"`
	Website 	string `form:"website,omitempty"`
	Phone		string `form:"phone,omitempty"`
	Email		string `form:"email,omitempty"`
	CompanySize	string `form:"company_size,omitempty"`
	Address		string `form:"address,omitempty"`
	Force		bool `form:"force,omitempty"`
	// NewContact	*Contact `form:"new_contact,omitempty"`
}