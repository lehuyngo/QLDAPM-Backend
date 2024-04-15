package entities

import "time"

type ClientAttachFile struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	UUID           string    `gorm:"column:uuid"`
	OriginalName   string    `gorm:"column:original_name;size:200"`
	RelativePath   string    `gorm:"column:relative_path_file;size:200"`
	Ext            string    `gorm:"column:ext"`
	OrganizationID int64     `gorm:"column:organization_id;omitempty"`
	ClientID       int64     `gorm:"column:client_id"`
	Client         *Client   `gorm:"references:client_id;foreignKey:id;omitempty"`
	CreatedBy      int64     `gorm:"column:created_by;omitempty"`
	Creator        *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;omitempty"`
	FileSize       int64     `gorm:"column:file_size;omitempty"`
}

func (*ClientAttachFile) TableName() string {
	return "client_attach_files"
}

func (c *ClientAttachFile) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *ClientAttachFile) GetOriginalName() string {
	if c == nil {
		return ""
	}

	return c.OriginalName
}

func (c *ClientAttachFile) GetRelativePath() string {
	if c == nil {
		return ""
	}

	return c.RelativePath
}

func (c *ClientAttachFile) GetExt() string {
	if c == nil {
		return ""
	}

	return c.Ext
}

func (c *ClientAttachFile) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}
