package entities

import "time"

type ContactMailActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	Contact   *Contact     `gorm:"references:contact_id;foreignKey:id;omitempty"`
	MailID    int64        `gorm:"column:mail_id;omitempty"`
	Mail      *Mail        `gorm:"references:mail_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (*ContactMailActivity) TableName() string {
	return "contact_mail_activities"
}

func (c *ContactMailActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ContactMailActivity) GetMail() *Mail {
	if c == nil {
		return nil
	}

	return c.Mail
}

func (c *ContactMailActivity) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}
