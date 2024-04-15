package entities

type MeetingNote struct {
	ID           int64                `gorm:"column:id;primaryKey"`
	ProjectID    int64                `gorm:"column:project_id;omitempty"`
	Project      *Project             `gorm:"references:project_id;foreignKey:id;omitempty"`
	MeetingID    int64                `gorm:"column:meeting_id;omitempty"`
	Meeting      *Meeting             `gorm:"references:meeting_id;foreignKey:id;omitempty"`
	StartTime    int64                `gorm:"column:start_time;omitempty"`
	Location     string               `gorm:"column:location;size:200;omitempty"`
	Link         string               `gorm:"column:link;size:200;omitempty"`
	Contributors []*Contributor       `gorm:"references:id;foreignKey:meeting_note_id"`
	Editors      []*MeetingNoteEditor `gorm:"references:id;foreignKey:meeting_note_id"`
	Highlights   []*MeetingHighlight  `gorm:"references:id;foreignKey:meeting_note_id"`
	Note         string               `gorm:"column:note;size:8192;omitempty"`
	*Base
}

func (*MeetingNote) TableName() string {
	return "meeting_notes"
}

func (c *MeetingNote) GetContributors() []*Contributor {
	if c == nil {
		return nil
	}

	return c.Contributors
}

func (c *MeetingNote) GetEditors() []*MeetingNoteEditor {
	if c == nil {
		return nil
	}

	return c.Editors
}

func (c *MeetingNote) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *MeetingNote) GetProject() *Project {
	if c == nil {
		return nil
	}

	return c.Project
}
