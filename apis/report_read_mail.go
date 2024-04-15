package apis

type ReadMailRecord struct {
	Sender 		*User `json:"sender,omitempty"`
	Receiver 	*Contact `json:"receiver,omitempty"`
	Mail 		*Mail `json:"mail,omitempty"`
	BatchMail 	*BatchMail `json:"batch_mail,omitempty"`
	SendTime 	int64 `json:"send_time,omitempty"`
	ReadTime 	int64 `json:"read_time,omitempty"`
	URL 		string `json:"url,omitempty"`
}

type ReportTimelineReadMailRequest struct {
	TimeRanges	[]*TimeRange `json:"time_ranges,omitempty"`
}

type ReportTimelineReadMailData struct {
	TimeRange	*TimeRange `json:"time_range,omitempty"`
	Records		[]*ReadMailRecord `json:"records,omitempty"`
}

type ReportTimelineReadMailResponse struct {
	Data []*ReportTimelineReadMailData `json:"data,omitempty"`
}
