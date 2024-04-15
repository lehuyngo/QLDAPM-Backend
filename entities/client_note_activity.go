package entities

import "time"

type ClientNoteActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	UUID      string       `gorm:"column:uuid;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	ClientID  int64        `gorm:"column:client_id;omitempty"`
	NoteID    int64        `gorm:"column:note_id;omitempty"`
	Note      *ClientNote  `gorm:"references:note_id;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
	Type      ActivityType `gorm:"column:type;omitempty"`
}

func (*ClientNoteActivity) TableName() string {
	return "client_note_activities"
}

func (c *ClientNoteActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ClientNoteActivity) GetNote() *ClientNote {
	if c == nil {
		return nil
	}

	return c.Note
}
