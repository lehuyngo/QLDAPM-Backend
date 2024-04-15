package entities

type ContactNote struct {
	ID        int64    `gorm:"column:id;primaryKey"`
	ContactID int64    `gorm:"column:contact_id;omitempty"`
	Contact   *Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	Title     string   `gorm:"column:title;size:200"`
	Content   string   `gorm:"column:content;size:4096"`
	Color     string   `gorm:"column:color;size:45"`
	Status    Status   `gorm:"column:status"`
	Base
}

func (*ContactNote) TableName() string {
	return "contact_notes"
}

func (r *ContactNote) GetID() int64 {
	if r == nil {
		return 0
	}

	return r.ID
}

func (c *ContactNote) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}

func (r *ContactNote) GetTitle() string {
	if r == nil {
		return ""
	}

	return r.Title
}

func (r *ContactNote) GetContent() string {
	if r == nil {
		return ""
	}

	return r.Content
}

func (r *ContactNote) GetColor() string {
	if r == nil {
		return ""
	}

	return r.Color
}
