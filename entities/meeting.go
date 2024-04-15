package entities

type Meeting struct {
	ID             int64       `gorm:"column:id;primaryKey"`
	ProjectID      int64       `gorm:"column:project_id;omitempty"`
	Project        *Project    `gorm:"references:project_id;foreignKey:id;omitempty"`
	StartTime      int64       `gorm:"column:start_time;omitempty"`
	Location       string      `gorm:"column:location;size:200;omitempty"`
	Link           string      `gorm:"column:link;size:200;omitempty"`
	Attendees      []*Attendee `gorm:"references:id;foreignKey:meeting_id"`
	LastActiveTime int64       `gorm:"column:last_active_time;omitempty"`
	*Base
}

func (*Meeting) TableName() string {
	return "meetings"
}

func (c *Meeting) GetAttendees() []*Attendee {
	if c == nil {
		return nil
	}

	return c.Attendees
}

func (c *Meeting) GetLastActiveTime() int64 {
	if c == nil {
		return 0
	}

	return c.LastActiveTime
}
