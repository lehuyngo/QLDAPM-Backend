package entities

import "time"

type ContactClientActivity struct {
	ID        int64        `gorm:"column:id;primaryKey"`
	Type      ActivityType `gorm:"column:type;omitempty"`
	ContactID int64        `gorm:"column:contact_id;omitempty"`
	Contact   *Contact     `gorm:"references:contact_id;foreignKey:id;omitempty"`
	ClientID  int64        `gorm:"column:client_id;omitempty"`
	Client    *Client      `gorm:"references:client_id;foreignKey:id;omitempty"`
	CreatedBy int64        `gorm:"column:created_by;omitempty"`
	Creator   *User        `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;omitempty"`
}
func (c *ContactClientActivity) TableName() string {		
	return "contact_client_activities"
}
func (c *ContactClientActivity) GetCreator() *User {
	if c == nil {
		return nil
	}

	return c.Creator
}

func (c *ContactClientActivity) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}
func (c *ContactClientActivity) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}	