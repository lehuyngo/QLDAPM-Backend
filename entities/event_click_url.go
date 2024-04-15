package entities

type EventClickURL struct {
	ID					int64 `gorm:"column:id;primaryKey"`
	Code				string `gorm:"column:code;omitempty"`
	SenderID			int64 `gorm:"column:sender_id;omitempty"`
	Sender   			*User `gorm:"references:sender_id;foreignKey:id;omitempty"`
	ReceiverID			int64 `gorm:"column:receiver_id;omitempty"`
	Receiver   			*Contact `gorm:"references:receiver_id;foreignKey:id;omitempty"`
	SendTime			int64 `gorm:"column:send_time;omitempty"`
	ReadTime			int64 `gorm:"column:read_time;omitempty"`
	URL					string `gorm:"column:url;omitempty"`
	OriginalURL			string `gorm:"column:original_url;omitempty"`
	MailID    			int64 `gorm:"column:mail_id;omitempty"`
	MailUUID    		string `gorm:"column:mail_uuid;omitempty"`
	MailSubject			string `gorm:"column:mail_subject;omitempty"`
	BatchMailID  		int64 `gorm:"column:batch_mail_id;omitempty"`
	BatchMailUUID  		string `gorm:"column:batch_mail_uuid;omitempty"`
	BatchMailSubject	string `gorm:"column:batch_mail_subject;omitempty"`
	OrganizationID		int64 `gorm:"column:organization_id;omitempty"`
}

func (*EventClickURL) TableName() string {
	return "event_click_urls"
}

func (c *EventClickURL) GetSender() *User {
	if c == nil {
		return nil
	}

	return c.Sender
}

func (c *EventClickURL) GetReceiver() *Contact {
	if c == nil {
		return nil
	}

	return c.Receiver
}
