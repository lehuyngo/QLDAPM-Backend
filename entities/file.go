package entities

import "time"

type File struct {
	UUID   				string `gorm:"column:uuid;primaryKey"`
	OriginalName 		string `gorm:"column:original_name;size:200"`
	RelativePath		string `gorm:"column:relative_path_file;size:200"`
	RelativeThumbnail	string `gorm:"column:relative_thumbnail;size:200"`
	Ext 				string `gorm:"column:ext"`
	CreatedBy 			int64 `gorm:"column:created_by;omitempty"`
	Creator   			*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 			time.Time `gorm:"column:created_at;omitempty"`
}

func (*File) TableName() string {
	return "files"
}

func (c *File) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *File) GetOriginalName() string {
	if c == nil {
		return ""
	}

	return c.OriginalName
}

func (c *File) GetRelativePath() string {
	if c == nil {
		return ""
	}

	return c.RelativePath
}

func (c *File) GetRelativeThumbnail() string {
	if c == nil {
		return ""
	}

	return c.RelativeThumbnail
}

func (c *File) GetExt() string {
	if c == nil {
		return ""
	}

	return c.Ext
}