package apis

type MiniProject struct {
	FullName 		string `json:"fullname,omitempty"`
	ShortName 		string `json:"shortname,omitempty"`
	Code			string `json:"code,omitempty"`
}

type CreateClientProjectRequest struct {
	UUIDs		[]string `json:"uuids,omitempty"`
	NewProject	*MiniProject `json:"new_project,omitempty"`
}
