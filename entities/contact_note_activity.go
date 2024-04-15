package entities

import "time"

type ContactNoteActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	NoteID    int64        `gorm:"column:note_id;omitempty"`
	Note      *ContactNote `gorm:"references:note_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (*ContactNoteActivity) TableName() string {
	return "contact_note_activities"
}

func (c *ContactNoteActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ContactNoteActivity) GetNote() *ContactNote {
	if c == nil {
		return nil
	}

	return c.Note
}
