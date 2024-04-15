package entities

import "time"

type ContactAttachFile struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	UUID   			string `gorm:"column:uuid"`
	OriginalName	string `gorm:"column:original_name;size:200"`
	RelativePath	string `gorm:"column:relative_path_file;size:200"`
	Ext 			string `gorm:"column:ext"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	ContactID 		int64 `gorm:"column:contact_id"`
	Contact 		*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	CreatedBy 		int64 `gorm:"column:created_by;omitempty"`
	Creator   		*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt		time.Time `gorm:"column:created_at;omitempty"`
}

func (*ContactAttachFile) TableName() string {
	return "contact_attach_files"
}

func (c *ContactAttachFile) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *ContactAttachFile) GetOriginalName() string {
	if c == nil {
		return ""
	}

	return c.OriginalName
}

func (c *ContactAttachFile) GetRelativePath() string {
	if c == nil {
		return ""
	}

	return c.RelativePath
}

func (c *ContactAttachFile) GetExt() string {
	if c == nil {
		return ""
	}

	return c.Ext
}

func (c *ContactAttachFile) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}
