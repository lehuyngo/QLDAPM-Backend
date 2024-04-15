package entities

import (
	"time"
)

type MeetingHighlight struct {
	ID            int64        `gorm:"column:id;primaryKey"`
	UUID          string       `gorm:"column:uuid;omitempty"`
	Title         string       `gorm:"column:title;size:8192"`
	MeetingNoteID int64        `gorm:"column:meeting_note_id;omitempty"`
	MeetingNote   *MeetingNote `gorm:"references:meeting_note_id;foreignKey:id;omitempty"`
	CreatedAt     time.Time    `gorm:"column:created_at;omitempty"`
	CreatedBy     int64        `gorm:"column:created_by;omitempty"`
	Creator       *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
}

func (*MeetingHighlight) TableName() string {
	return "meeting_highlights"
}

func (c *MeetingHighlight) GetMeetingNote() *MeetingNote {
	if c == nil {
		return nil
	}

	return c.MeetingNote
}
