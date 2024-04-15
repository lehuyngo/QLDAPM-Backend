package apis

type MiniClient struct {
	FullName 		string `json:"fullname,omitempty"`
	ShortName 		string `json:"shortname,omitempty"`
	Code 			string `json:"code,omitempty"`
	Fax 			string `json:"fax,omitempty"`
	Website 		string `json:"website,omitempty"`
	Phone			string `json:"phone,omitempty"`
	Email			string `json:"email,omitempty"`
	CompanySize		string `json:"company_size,omitempty"`
	Address			string `json:"address,omitempty"`
}

type CreateContactClientRequest struct {
	UUIDs		[]string `json:"uuids,omitempty"`
	NewClient	*MiniClient `json:"new_client,omitempty"`
}
