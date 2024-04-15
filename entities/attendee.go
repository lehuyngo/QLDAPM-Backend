package entities

import "time"

type Attendee struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	UUID      string    `gorm:"column:uuid;omitempty"`
	MeetingID int64     `gorm:"column:meeting_id;omitempty"`
	Meeting   *Meeting  `gorm:"references:meeting_id;foreignKey:id;omitempty"`
	ContactID int64     `gorm:"column:contact_id;omitempty"`
	Contact   *Contact  `gorm:"references:contact_id;foreignKey:id;omitempty"`
	UserID    int64     `gorm:"column:user_id;omitempty"`
	User      *User     `gorm:"references:user_id;foreignKey:id;omitempty"`
	CreatedBy int64     `gorm:"column:created_by;omitempty"`
	Creator   *User     `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;omitempty"`
}

func (*Attendee) TableName() string {
	return "attendees"
}

func (c *Attendee) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *Attendee) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}

func (c *Attendee) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}

func (c *Attendee) GetUser() *User {
	if c == nil {
		return nil
	}

	return c.User
}

func (c *Attendee) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}
