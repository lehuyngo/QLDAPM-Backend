package apis

type ContactNote struct {
	UUID   		string `json:"uuid,omitempty"`
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
	CreateTime	int64 `json:"create_time,omitempty"`
	Creator 	*User `json:"creator,omitempty"`
}

type CreateContactNoteRequest struct {
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
}

type UpdateContactNoteRequest struct {
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
}

type ListContactNoteResponse struct {
	Data []*ContactNote `json:"data"`
}