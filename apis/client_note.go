package apis

type ClientNote struct {
	UUID   		string `json:"uuid,omitempty"`
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
	CreateTime	int64 `json:"create_time,omitempty"`
	Creator 	*User `json:"creator,omitempty"`
}

type CreateClientNoteRequest struct {
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
}

type UpdateClientNoteRequest struct {
	Title		string `json:"title,omitempty"`
	Content 	string `json:"content,omitempty"`
	Color		string `json:"color,omitempty"`
}

type ListClientNoteResponse struct {
	Data []*ClientNote `json:"data"`
}