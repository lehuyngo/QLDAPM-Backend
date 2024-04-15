package entities

import "time"

type TaskComment struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	UUID      string    `gorm:"column:uuid;omitempty"`
	Content   string    `gorm:"column:content;size:4096"`
	TaskID    int64     `gorm:"column:task_id;omitempty"`
	Task      *Task     `gorm:"references:task_id;foreignKey:id;omitempty"`
	CreatedBy int64     `gorm:"column:created_by;omitempty"`
	Creator   *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;omitempty"`
}

func (*TaskComment) TableName() string {
	return "task_comments"
}

func (c *TaskComment) GetID() int64 {
	if c == nil {
		return 0
	}
	return c.ID
}

func (c *TaskComment) GetTask() *Task {
	if c == nil {
		return nil
	}
	return c.Task
}

func (c *TaskComment) GetContent() string {
	if c == nil {
		return ""
	}
	return c.Content
}

func (c *TaskComment) GetCreator() *User {
	if c == nil {
		return nil
	}
	return c.Creator
}
