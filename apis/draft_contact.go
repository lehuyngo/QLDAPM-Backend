package apis

type DraftContact struct {
	UUID   			string `json:"uuid,omitempty"`
	FullName 		string `json:"fullname,omitempty"`
	Phone			string `json:"phone,omitempty"`
	Email			string `json:"email,omitempty"`
	ClientName 		string `json:"client_name,omitempty"`
	ClientWebsite	string `json:"client_website,omitempty"`
	ClientAddress	string `json:"client_address,omitempty"`
	NameCard		*File `json:"name_card,omitempty"`
	CompanyLogo		*File `json:"company_logo,omitempty"`
}

type CreateDraftContactRequest struct {
	FullName 		string `form:"fullname,omitempty"`
	Phone			string `form:"phone,omitempty"`
	Email			string `form:"email,omitempty"`
	ClientName 		string `form:"client_name,omitempty"`
	ClientWebsite	string `form:"client_website,omitempty"`
	ClientAddress	string `form:"client_address,omitempty"`
}

type ListDraftContactResponse struct {
	Data []*DraftContact `json:"data"`
}

type ConvertDraftContactRequest struct {
	FullName 		string `form:"fullname,omitempty"`
	Phone			string `form:"phone,omitempty"`
	Email			string `form:"email,omitempty"`
	Tags			string `form:"tags,omitempty"`
	ClientName 		string `form:"client_name,omitempty"`
	ClientWebsite	string `form:"client_website,omitempty"`
	ClientAddress	string `form:"client_address,omitempty"`
	ClientTags		string `form:"client_tags,omitempty"`
}

type ConvertDraftContactResponse struct {
	ContactUUID string `json:"contact_uuid"`
	ClientUUID string `json:"client_uuid"`
}