package apis

type MiniContact struct {
	FullName 		string `json:"fullname,omitempty"`
	ShortName 		string `json:"shortname,omitempty"`
	Phone			string `json:"phone,omitempty"`
	Email			string `json:"email,omitempty"`
	JobTitle		string `json:"job_title,omitempty"`
	Gender			int32 `json:"gender,omitempty"`
	LastActiveTime	int64 `json:"last_active_time,omitempty"`
	BirthDay		int64 `json:"birthday,omitempty"`
}

type CreateClientContactRequest struct {
	UUIDs		[]string `json:"uuids,omitempty"`
	NewContact	*MiniContact `json:"new_contact,omitempty"`
}
