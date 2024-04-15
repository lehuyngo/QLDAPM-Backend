package apis

type CreateContactProjectRequest struct {
	UUIDs		[]string `json:"uuids,omitempty"`
	NewProject	*MiniProject `json:"new_project,omitempty"`
}
