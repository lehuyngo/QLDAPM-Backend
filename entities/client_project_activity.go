package entities

import "time"

type ClientProjectActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ClientID  int64        `gorm:"column:client_id;omitempty"`
	Client    *Client      `gorm:"references:client_id;foreignKey:id;omitempty"`
	ProjectID int64        `gorm:"column:project_id;omitempty"`
	Project   *Project     `gorm:"references:project_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (c *ClientProjectActivity) TableName() string {
	return "client_project_activities"
}
func (c *ClientProjectActivity) GetCreator() *User {
	if c == nil {
		return nil
	}
	return c.Creator
}
func (c *ClientProjectActivity) GetClient() *Client {
	if c == nil {
		return nil
	}
	return c.Client
}
func (c *ClientProjectActivity) GetProject() *Project {
	if c == nil {
		return nil
	}
	return c.Project
}
