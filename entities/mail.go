package entities

import "time"

type Mail struct {
	ID             	int64 				`gorm:"column:id;primaryKey"`
	Subject       	string 				`gorm:"column:subject;size:1024"`
	Content     	string 				`gorm:"column:content;size:4096"`
	OrganizationID	int64 				`gorm:"column:organization_id;omitempty"`
	Receivers		[]*MailReceiver 	`gorm:"references:id;foreignKey:mail_id"`
	CarbonCopies   	[]*MailCarbonCopy 	`gorm:"references:id;foreignKey:mail_id"`
	AttachFiles    	[]*MailAttachFile 	`gorm:"references:id;foreignKey:mail_id"`
	URLs			[]*TrackedURL		`gorm:"references:id;foreignKey:mail_id"`
	Base
}

func (*Mail) TableName() string {
	return "mails"
}

type MailReceiver struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	UUID        string    `gorm:"column:uuid;omitempty"`
	MailID      int64     `gorm:"column:mail_id;omitempty"`
	Mail        *Mail     `gorm:"references:mail_id;foreignKey:id;omitempty"`
	ContactID   int64     `gorm:"column:contact_id;omitempty"`
	Contact     *Contact  `gorm:"references:contact_id;foreignKey:id;omitempty"`
	UserID      int64     `gorm:"column:user_id;omitempty"`
	User        *User     `gorm:"references:user_id;foreignKey:id;omitempty"`
	MailAddress string    `gorm:"column:mail_address;omitempty"`
	CreatedBy   int64     `gorm:"column:created_by;omitempty"`
	Creator     *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;omitempty"`
}

func (*MailReceiver) TableName() string {
	return "mail_receivers"
}

type MailCarbonCopy struct {
	ID   		int64 `gorm:"column:id;primaryKey"`
	UUID		string `gorm:"column:mail_id;omitempty"`
	MailID		int64 `gorm:"column:mail_id;omitempty"`
	Mail		*Mail `gorm:"references:mail_id;foreignKey:id;omitempty"`
	ContactID	int64 `gorm:"column:contact_id;omitempty"`
	Contact 	*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	UserID		int64 `gorm:"column:user_id;omitempty"`
	User 		*User `gorm:"references:user_id;foreignKey:id;omitempty"`
	MailAddress	string `gorm:"column:mail_address;omitempty"`
	CreatedBy 	int64 `gorm:"column:created_by;omitempty"`
	Creator   	*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 	time.Time `gorm:"column:created_at;omitempty"`
}

func (*MailCarbonCopy) TableName() string {
	return "mail_carbon_copies"
}

func (c *Mail) GetSubject() string {
	if c == nil {
		return ""
	}

	return c.Subject
}
