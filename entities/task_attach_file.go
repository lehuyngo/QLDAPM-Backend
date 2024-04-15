package entities

import "time"

type TaskAttachFile struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	UUID           string    `gorm:"column:uuid"`
	OriginalName   string    `gorm:"column:original_name;size:200"`
	RelativePath   string    `gorm:"column:relative_path_file;size:200"`
	Ext            string    `gorm:"column:ext"`
	OrganizationID int64     `gorm:"column:organization_id;omitempty"`
	TaskID         int64     `gorm:"column:task_id"`
	Task           *Task     `gorm:"references:task_id;foreignKey:id;omitempty"`
	CreatedBy      int64     `gorm:"column:created_by;omitempty"`
	Creator        *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at;omitempty"`
}

func (*TaskAttachFile) TableName() string {
	return "task_attach_files"
}

func (c *TaskAttachFile) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *TaskAttachFile) GetOriginalName() string {
	if c == nil {
		return ""
	}

	return c.OriginalName
}

func (c *TaskAttachFile) GetRelativePath() string {
	if c == nil {
		return ""
	}

	return c.RelativePath
}

func (c *TaskAttachFile) GetExt() string {
	if c == nil {
		return ""
	}

	return c.Ext
}

func (c *TaskAttachFile) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}

func (c *TaskAttachFile) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}
