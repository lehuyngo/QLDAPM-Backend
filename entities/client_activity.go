package entities

import "time"

type ClientActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
	ClientID  int64        `gorm:"column:client_id;omitempty"`
	Client    *Client      `gorm:"references:client_id;foreignKey:id;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
}

func (c *ClientActivity) TableName() string {
	return "client_activities"
}
func (c *ClientActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}
func (c *ClientActivity) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}
