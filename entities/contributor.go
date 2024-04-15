package entities

import "time"

type Contributor struct {
	ID            int64        `gorm:"column:id;primaryKey"`
	UUID          string       `gorm:"column:uuid;omitempty"`
	MeetingNoteID int64        `gorm:"column:meeting_note_id;omitempty"`
	MeetingNote   *MeetingNote `gorm:"references:meeting_note_id;foreignKey:id;omitempty"`
	ContactID     int64        `gorm:"column:contact_id;omitempty"`
	Contact       *Contact     `gorm:"references:contact_id;foreignKey:id;omitempty"`
	UserID        int64        `gorm:"column:user_id;omitempty"`
	User          *User        `gorm:"references:user_id;foreignKey:id;omitempty"`
	CreatedBy     int64        `gorm:"column:created_by;omitempty"`
	Creator       *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt     time.Time    `gorm:"column:created_at;omitempty"`
}

func (*Contributor) TableName() string {
	return "contributors"
}

func (c *Contributor) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *Contributor) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}

func (c *Contributor) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}

func (c *Contributor) GetUser() *User {
	if c == nil {
		return nil
	}

	return c.User
}

func (c *Contributor) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *Contributor) GetMeetingNote() *MeetingNote {
	if c == nil {
		return nil
	}

	return c.MeetingNote
}
