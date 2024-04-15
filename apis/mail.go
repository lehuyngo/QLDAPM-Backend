package apis

type CreateMailRequest struct {
	Subject					string `form:"subject,omitempty"`
	Content					string `form:"content,omitempty"`
	ReceiverContactUUIDs	string `form:"receiver_contact_uuids,omitempty"`
	ReceiverUserUUIDs		string `form:"receiver_user_uuids,omitempty"`
	ReceiverMailAddresses	string `form:"receiver_mail_addresses,omitempty"`
	CCContactUUIDs			string `form:"cc_contact_uuids,omitempty"`
	CCUserUUIDs				string `form:"cc_user_uuids,omitempty"`
	CCMailAddresses			string `form:"cc_mail_addresses,omitempty"`
	UseShortenLink			int32 `form:"use_shorten_link,omitempty"`
}

type Mail struct {
	UUID   					string `json:"uuid,omitempty"`
	Sender					*User `json:"sender,omitempty"`
	SendTime				int64 `json:"send_time,omitempty"`
	Subject 				string `json:"subject,omitempty"`
	Content 				string `json:"content,omitempty"`
	ReceiverContacts		[]*Contact `json:"receiver_contacts,omitempty"`
	ReceiverUsers			[]*User `json:"receiver_users,omitempty"`
	ReceiverMailAddresses	[]string `json:"receiver_mail_addresses,omitempty"`
	CCContacts				[]*Contact `json:"cc_contacts,omitempty"`
	CCUsers					[]*User `json:"cc_users,omitempty"`
	CCMailAddresses			[]string `json:"cc_mail_addresses,omitempty"`
}

type ListMailResponse struct {
	Data []*Mail `json:"data"`
}

type ReadMailRequest struct {
	Code string `json:"code,omitempty"`
}
