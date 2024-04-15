package entities

import "time"

type MailAttachFile struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	UUID   			string `gorm:"column:uuid"`
	OriginalName	string `gorm:"column:original_name;size:200"`
	RelativePath	string `gorm:"column:relative_path_file;size:200"`
	Ext 			string `gorm:"column:ext"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	MailID 			int64 `gorm:"column:mail_id"`
	Mail 			*Mail `gorm:"references:mail_id;foreignKey:id;omitempty"`
	CreatedBy 		int64 `gorm:"column:created_by;omitempty"`
	Creator   		*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt		time.Time `gorm:"column:created_at;omitempty"`
}

func (*MailAttachFile) TableName() string {
	return "mail_attach_files"
}

func (c *MailAttachFile) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *MailAttachFile) GetOriginalName() string {
	if c == nil {
		return ""
	}

	return c.OriginalName
}

func (c *MailAttachFile) GetRelativePath() string {
	if c == nil {
		return ""
	}

	return c.RelativePath
}

func (c *MailAttachFile) GetExt() string {
	if c == nil {
		return ""
	}

	return c.Ext
}

func (c *MailAttachFile) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}
