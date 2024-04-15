package entities

import "time"

type ContactMailShortClick struct {
	ID   int64  `gorm:"column:id;primaryKey"`
	UUID string `gorm:"column:uuid;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	Contact   *Contact     `gorm:"references:contact_id;foreignKey:id;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;omitempty"`
}

func (c *ContactMailShortClick) TableName() string {
	return "contact_mail_shortclicks"
}
func (c *ContactMailShortClick) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}
func (c *ContactMailShortClick) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}
