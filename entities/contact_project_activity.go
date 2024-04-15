package entities

import "time"

type ContactProjectActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	Contact   *Contact     `gorm:"references:contact_id;foreignKey:id;omitempty"`
	ProjectID int64        `gorm:"column:project_id;omitempty"`
	Project   *Project     `gorm:"references:project_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (c *ContactProjectActivity) TableName() string {
	return "contact_project_activities"
}
func (c *ContactProjectActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ContactProjectActivity) GetProject() *Project {
	if c == nil {
		return nil
	}

	return c.Project
}
