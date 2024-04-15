package apis

type ReceiveMailHistory struct {
	UUID      		string `json:"uuid,omitempty"`
	Code      		string `json:"code,omitempty"`
	URL      		string `json:"url,omitempty"`
	Mail      		*Mail `json:"mail,omitempty"`
	BatchMail 		*BatchMail `json:"batch_mail,omitempty"`
	Creator  		*User `json:"creator,omitempty"`
	CreatedTime		int64 `json:"created_time,omitempty"`
	ReceivedTime	int64 `json:"received_time,omitempty"`
	Status 			int32 `json:"status,omitempty"`
}

type ListReceiveMailHistoryResponse struct {
	Data []*ReceiveMailHistory `json:"data"`
}