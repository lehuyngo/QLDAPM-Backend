package entities

import "time"

type ContactTagActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	TagID     int64        `gorm:"column:tag_id;omitempty"`
	Tag       *ContactTag  `gorm:"references:tag_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (c *ContactTagActivity) TableName() string {
	return "contact_tag_activities"
}

func (c *ContactTagActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ContactTagActivity) GetTag() *ContactTag {
	if c == nil {
		return nil
	}

	return c.Tag
}
