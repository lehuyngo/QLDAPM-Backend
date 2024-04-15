package apis

type CreateBatchMailRequest struct {
	Subject					string `form:"subject,omitempty"`
	Content					string `form:"content,omitempty"`
	ReceiverContactUUIDs	string `form:"receiver_contact_uuids,omitempty"`
	CCUserUUIDs				string `form:"cc_user_uuids,omitempty"`
	CCMailAddresses			string `form:"cc_mail_addresses,omitempty"`
	UseShortenLink			int32 `form:"use_shorten_link,omitempty"`
}

type BatchMail struct {
	UUID   			string `json:"uuid,omitempty"`
	Sender			*User `json:"sender,omitempty"`
	SendTime		int64 `json:"send_time,omitempty"`
	Subject 		string `json:"subject,omitempty"`
	Content 		string `json:"content,omitempty"`
	Status 			int `json:"status,omitempty"`
	Receivers		[]*BatchMailReceiver `json:"receivers,omitempty"`
	CCUsers			[]*User `json:"cc_users,omitempty"`
	CCMailAddresses	[]string `json:"cc_mail_addresses,omitempty"`
}

type BatchMailReceiver struct {
	UUID		string `json:"uuid,omitempty"`
	Contact		*Contact `json:"contact,omitempty"`
	SendTime	int64 `json:"send_time,omitempty"`
	Status		int `json:"status,omitempty"`
}

type BatchMailURL struct {
	UUID	string `json:"uuid,omitempty"`
	Alias	string `json:"alias,omitempty"`
	URL		string `json:"url,omitempty"`
}

type ListBatchMailResponse struct {
	Data []*BatchMail `json:"data"`
}
