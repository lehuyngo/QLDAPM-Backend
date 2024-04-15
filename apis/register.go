package apis

type Register struct {
	DisplayName      string `json:"displayname,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	Email            string `json:"email,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
}

type User struct {
	UUID        string `json:"uuid,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
}

type SelfProfile struct {
	UUID        string `json:"uuid,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
}

type ListSelfProfile struct {
	Data []*SelfProfile `json:"data,omitempty"`
}
