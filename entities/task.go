package entities

type Task struct {
	ID             int64             `gorm:"column:id;primaryKey"`
	Title          string            `gorm:"column:title;size:255;omitempty"`
	Status         TaskStatus        `gorm:"column:status;omitempty"`
	Priority       TaskPriority      `gorm:"column:priority;omitempty"`
	Label          TaskLabel         `gorm:"column:label;omitempty"`
	ProjectID      int64             `gorm:"column:project_id;omitempty"`
	Project        *Project          `gorm:"references:project_id;foreignKey:id;omitempty"`
	Deadline       int64             `gorm:"column:deadline;omitempty"`
	EstimatedHours int32             `gorm:"column:estimated_hours;omitempty"`
	Description    string            `gorm:"column:description;size:4096"`
	Assignees      []*TaskAssignee   `gorm:"references:id;foreignKey:task_id"`
	OrganizationID int64             `gorm:"column:organization_id;omitempty"`
	AttachFiles    []*TaskAttachFile `gorm:"references:id;foreignKey:task_id"`
	Base
}

func (*Task) TableName() string {
	return "tasks"
}

func (c *Task) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *Task) GetTitle() string {
	if c == nil {
		return ""
	}

	return c.Title
}

func (c *Task) GetProject() *Project {
	if c == nil {
		return nil
	}

	return c.Project
}

func (c *Task) GetAssignees() []*TaskAssignee {
	if c == nil {
		return nil
	}

	return c.Assignees
}

func (c *Task) GetAttachFiles() []*TaskAttachFile {
	if c == nil {
		return nil
	}

	return c.AttachFiles
}
