package entities

import "time"

type MeetingNoteEditor struct {
	ID            int64        `gorm:"column:id;primaryKey"`
	UUID          string       `gorm:"column:uuid;omitempty"`
	MeetingNoteID int64        `gorm:"column:meeting_note_id;omitempty"`
	MeetingNote   *MeetingNote `gorm:"references:meeting_note_id;foreignKey:id;omitempty"`
	EditorID      int64        `gorm:"column:editor_id;omitempty"`
	Editor        *User        `gorm:"references:editor_id;foreignKey:id;omitempty"`
	CreatedBy     int64        `gorm:"column:created_by;omitempty"`
	Creator       *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt     time.Time    `gorm:"column:created_at;omitempty"`
}

func (*MeetingNoteEditor) TableName() string {
	return "meeting_note_editors"
}

func (c *MeetingNoteEditor) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

func (c *MeetingNoteEditor) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}

func (c *MeetingNoteEditor) GetEditor() *User {
	if c == nil {
		return nil
	}

	return c.Editor
}

func (c *MeetingNoteEditor) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}
