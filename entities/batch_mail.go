package entities

import "time"

type BatchMail struct {
	ID    			int64 `gorm:"column:id;primaryKey"`
	Subject     	string `gorm:"column:subject;size:1024"`
	Content      	string `gorm:"column:content;size:4096"`
	OrganizationID 	int64 `gorm:"column:organization_id;omitempty"`
	Status 			MailStatus `gorm:"column:status;omitempty"`
	UserShortenLink	int32 `gorm:"column:user_shorten_link;omitempty"`
	Receivers		[]*BatchMailReceiver `gorm:"references:id;foreignKey:mail_id"`
	CarbonCopies  	[]*BatchMailCarbonCopy `gorm:"references:id;foreignKey:mail_id"`
	AttachFiles    	[]*BatchMailAttachFile `gorm:"references:id;foreignKey:mail_id"`
	URLs			[]*TrackedURL `gorm:"references:id;foreignKey:batch_mail_id"`
	Base
}

func (*BatchMail) TableName() string {
	return "batch_mails"
}

func (c *BatchMail) GetSubject() string {
	if c == nil {
		return ""
	}

	return c.Subject
}

type BatchMailReceiver struct {
	ID          int64     	`gorm:"column:id;primaryKey"`
	UUID        string    	`gorm:"column:uuid;omitempty"`
	MailID      int64     	`gorm:"column:mail_id;omitempty"`
	Mail        *BatchMail	`gorm:"references:mail_id;foreignKey:id;omitempty"`
	Email   	string  	`gorm:"column:email;omitempty"`
	ContactID   int64     	`gorm:"column:contact_id;omitempty"`
	Contact     *Contact  	`gorm:"references:contact_id;foreignKey:id;omitempty"`
	Status 		MailStatus 	`gorm:"column:status;omitempty"`
	SendTime  	int64     	`gorm:"column:send_time;omitempty"`
	CreatedBy   int64     	`gorm:"column:created_by;omitempty"`
	Creator     *User     	`gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt   time.Time 	`gorm:"column:created_at;omitempty"`
}

func (*BatchMailReceiver) TableName() string {
	return "batch_mail_receivers"
}

func (b *BatchMailReceiver) GetContact() *Contact {
	if b == nil {
		return nil
	}

	return b.Contact
}

type BatchMailCarbonCopy struct {
	ID   		int64 `gorm:"column:id;primaryKey"`
	UUID		string `gorm:"column:mail_id;omitempty"`
	MailID		int64 `gorm:"column:mail_id;omitempty"`
	Mail		*BatchMail `gorm:"references:mail_id;foreignKey:id;omitempty"`
	UserID		int64 `gorm:"column:user_id;omitempty"`
	User 		*User `gorm:"references:user_id;foreignKey:id;omitempty"`
	MailAddress	string `gorm:"column:mail_address;omitempty"`
	CreatedBy 	int64 `gorm:"column:created_by;omitempty"`
	Creator   	*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 	time.Time `gorm:"column:created_at;omitempty"`
}

func (*BatchMailCarbonCopy) TableName() string {
	return "batch_mail_carbon_copies"
}

type BatchMailAttachFile struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	UUID   			string `gorm:"column:uuid"`
	OriginalName	string `gorm:"column:original_name;size:200"`
	RelativePath	string `gorm:"column:relative_path_file;size:200"`
	Ext 			string `gorm:"column:ext"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	MailID 			int64 `gorm:"column:mail_id"`
	Mail 			*BatchMail `gorm:"references:mail_id;foreignKey:id;omitempty"`
	CreatedBy 		int64 `gorm:"column:created_by;omitempty"`
	Creator   		*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt		time.Time `gorm:"column:created_at;omitempty"`
}

func (*BatchMailAttachFile) TableName() string {
	return "batch_mail_attach_files"
}
