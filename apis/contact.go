package apis

type Contact struct {
	UUID   			string `json:"uuid,omitempty"`
	FullName 		string `json:"fullname,omitempty"`
	ShortName 		string `json:"shortname,omitempty"`
	Phone			string `json:"phone,omitempty"`
	Email			string `json:"email,omitempty"`
	JobTitle		string `json:"job_title,omitempty"`
	Gender			int32 `json:"gender,omitempty"`
	LastActiveTime	int64 `json:"last_active_time,omitempty"`
	BirthDay		int64 `json:"birthday,omitempty"`
	CreatedTime		int64 `json:"created_time,omitempty"`
	Avatar 			*File `json:"avatar,omitempty"`
	NameCard		*File `json:"name_card,omitempty"`
	Tags			[]*ContactTag `json:"tags,omitempty"`
	Notes			[]*ContactNote `json:"notes,omitempty"`
	Clients			[]*Client `json:"clients,omitempty"`
}

func (c *Contact) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

type ListContactResponse struct {
	Data []*Contact `json:"data"`
}

type CreateContactRequest struct {
	FullName 	string `form:"fullname,omitempty"`
	ShortName 	string `form:"shortname,omitempty"`
	Phone		string `form:"phone,omitempty"`
	Email		string `form:"email,omitempty"`
	JobTitle	string `form:"job_title,omitempty"`
	Gender		int32 `form:"gender,omitempty"`
	BirthDay	int64 `form:"birthday,omitempty"`
	ClientUUIDs []string `form:"client_uuids,omitempty"`
}
