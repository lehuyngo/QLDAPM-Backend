package entities

import "time"

type TrackedURL struct {
	ID        		int64 `gorm:"column:id;primaryKey"`
	UUID      		string `gorm:"column:uuid;omitempty"`
	Code      		string `gorm:"column:code;omitempty"`
	OriginalURL    	string `gorm:"column:original_url;omitempty"`
	URL      		string `gorm:"column:url;omitempty"`
	ContactID    	int64 `gorm:"column:contact_id;omitempty"`
	Contact      	*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	MailID    		int64 `gorm:"column:mail_id;omitempty"`
	Mail      		*Mail `gorm:"references:mail_id;foreignKey:id;omitempty"`
	BatchMailID		int64 `gorm:"column:batch_mail_id;omitempty"`
	BatchMail 		*BatchMail `gorm:"references:batch_mail_id;foreignKey:id;omitempty"`
	CreatedBy		int64 `gorm:"column:created_by;omitempty"`
	Creator  		*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt		time.Time `gorm:"column:created_at;omitempty"`
	ReceivedAt		int64 `gorm:"column:received_at;omitempty"`
	Status   		ReadStatus `gorm:"column:status;omitempty"`
	OrganizationID 	int64 `gorm:"column:organization_id;omitempty"`
}

func (*TrackedURL) TableName() string {
	return "tracked_urls"
}

func (c *TrackedURL) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *TrackedURL) GetMail() *Mail {
	if c == nil {
		return nil
	}

	return c.Mail
}

func (c *TrackedURL) GetBatchMail() *BatchMail {
	if c == nil {
		return nil
	}

	return c.BatchMail
}

func (c *TrackedURL) GetReceivedAt() int64 {
	if c == nil {
		return 0
	}

	return c.ReceivedAt
}

func (c *TrackedURL) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}
