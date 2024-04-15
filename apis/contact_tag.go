package apis

type ContactTag struct {
	UUID		string `json:"uuid,omitempty"`
	Name		string `json:"name,omitempty"`
	Color		string `json:"color,omitempty"`
	CreatedTime	int64 `json:"created_time,omitempty"`
	Contacts 	[]*Contact `json:"contacts,omitempty"`
}

type CreateContactTagRequest struct {
	Name	string `json:"name,omitempty"`
	Color	string `json:"color,omitempty"`
}

type ListContactTagRequest struct {
	Data []*ContactTag `json:"data,omitempty"`
}