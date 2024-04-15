package entities

import "time"

type ClientTagActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ClientID  int64        `gorm:"column:client_id;omitempty"`
	TagID     int64        `gorm:"column:tag_id;omitempty"`
	Tag       *ClientTag   `gorm:"references:tag_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (*ClientTagActivity) TableName() string {
	return "client_tag_activities"
}

func (c *ClientTagActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ClientTagActivity) GetTag() *ClientTag {
	if c == nil {
		return nil
	}

	return c.Tag
}
