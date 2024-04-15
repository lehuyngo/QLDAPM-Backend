package entities

import "time"

type ContactProject struct {
	ContactID 	int64 `gorm:"column:contact_id;primaryKey"`
	Contact 	*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	ProjectID 	int64 `gorm:"column:project_id;primaryKey"`
	Project 	*Project `gorm:"references:project_id;foreignKey:id;omitempty"`
	CreatedBy 	int64 `gorm:"column:created_by;omitempty"`
	Creator   	*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 	time.Time `gorm:"column:created_at;omitempty"`
}

func (*ContactProject) TableName() string {
	return "contact_projects"
}

func (c *ContactProject) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}

func (c *ContactProject) GetProject() *Project {
	if c == nil {
		return nil
	}

	return c.Project
}
