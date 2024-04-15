package entities

import "time"

type TaskAssignee struct {
	TaskID     int64     `gorm:"column:task_id;primaryKey"`
	Task       *Task     `gorm:"references:task_id;foreignKey:id;omitempty"`
	AssigneeID int64     `gorm:"column:assignee_id;primaryKey"`
	Assignee   *User     `gorm:"references:assignee_id;foreignKey:id;omitempty"`
	CreatedBy  int64     `gorm:"column:created_by;omitempty"`
	Creator    *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;omitempty"`
}

func (*TaskAssignee) TableName() string {
	return "task_assignees"
}

func (c *TaskAssignee) GetTask() *Task {
	if c == nil {
		return nil
	}

	return c.Task
}

func (c *TaskAssignee) GetAssignee() *User {
	if c == nil {
		return nil
	}

	return c.Assignee
}
