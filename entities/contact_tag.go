package entities

import (
	"time"

	"gorm.io/gorm"
)

type ContactTag struct {
	ID             int64                `gorm:"column:id;primaryKey"`
	UUID           string               `gorm:"column:uuid"`
	Name           string               `gorm:"column:name;size:200"`
	Color          string               `gorm:"column:color;size:45"`
	CreatedBy      int64                `gorm:"column:created_by;omitempty"`
	Creator        *User                `gorm:"references:created_by;foreignKey:id;omitempty"`
	OrganizationID int64                `gorm:"column:organization_id;omitempty"`
	CreatedAt      time.Time            `gorm:"column:created_at;omitempty"`
	DeletedAt      gorm.DeletedAt       `swaggertype:"string" gorm:"index;column:deleted_at"`
	Contacts       []*ContactContactTag `gorm:"references:id;foreignKey:tag_id"`
}

func (*ContactTag) TableName() string {
	return "contact_tags"
}

func (r *ContactTag) GetID() int64 {
	if r == nil {
		return 0
	}

	return r.ID
}

func (r *ContactTag) GetUUID() string {
	if r == nil {
		return ""
	}

	return r.UUID
}

func (r *ContactTag) GetName() string {
	if r == nil {
		return ""
	}

	return r.Name
}

func (r *ContactTag) GetColor() string {
	if r == nil {
		return "#FFFFFF"
	}

	if r.Color == "" {
		return "#FFFFFF"
	}

	return r.Color
}

func (r *ContactTag) GetCreatedAt() time.Time {
	if r == nil {
		return time.Now()
	}

	return r.CreatedAt
}

type ContactContactTag struct {
	ContactID int64       `gorm:"column:contact_id;primaryKey"`
	Contact   *Contact    `gorm:"references:contact_id;foreignKey:id;omitempty"`
	TagID     int64       `gorm:"column:tag_id;primaryKey"`
	Tag       *ContactTag `gorm:"references:tag_id;foreignKey:id;omitempty"`
	Color     string      `gorm:"column:color;size:45"`
	CreatedBy int64       `gorm:"column:created_by;omitempty"`
	Creator   *User       `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time   `gorm:"column:created_at;omitempty"`
}

func (*ContactContactTag) TableName() string {
	return "contact_contact_tags"
}

func (c *ContactContactTag) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}

func (c *ContactContactTag) GetTag() *ContactTag {
	if c == nil {
		return nil
	}

	return c.Tag
}

func (c *ContactContactTag) GetCreatedAt() time.Time {
	if c == nil {
		return time.Now()
	}

	return c.CreatedAt
}

func (c *ContactContactTag) GetColor() string {
	if c == nil {
		return ""
	}

	return c.Color
}