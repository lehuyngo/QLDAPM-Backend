package apis

type ClientTag struct {
	UUID        string    `json:"uuid,omitempty"`
	Name        string    `json:"name,omitempty"`
	Color       string    `json:"color,omitempty"`
	CreatedTime int64     `json:"created_time,omitempty"`
	Clients     []*Client `json:"clients,omitempty"`
}

type CreateClientTagRequest struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type ListClientTagRequest struct {
	Data []*ClientTag `json:"data,omitempty"`
}

type DeleteClientTagRequest struct {
	IsFullyDeleted int32 `json:"is_fully_deleted,omitempty"`
}
