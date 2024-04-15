package entities

import "time"

type ClientContact struct {
	ClientID 	int64 `gorm:"column:client_id;primaryKey"`
	Client 		*Client `gorm:"references:client_id;foreignKey:id;omitempty"`
	ContactID 	int64 `gorm:"column:contact_id;primaryKey"`
	Contact 	*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	CreatedBy 	int64 `gorm:"column:created_by;omitempty"`
	Creator   	*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 	time.Time `gorm:"column:created_at;omitempty"`
}

func (*ClientContact) TableName() string {
	return "client_contacts"
}

func (c *ClientContact) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}

func (c *ClientContact) GetContact() *Contact {
	if c == nil {
		return nil
	}

	return c.Contact
}