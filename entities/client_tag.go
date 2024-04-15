package entities

import (
	"time"

	"gorm.io/gorm"
)

type ClientTag struct {
	ID		int64 `gorm:"column:id;primaryKey"`
	UUID	string `gorm:"column:uuid"`
	Name	string `gorm:"column:name;size:200"`
	Color	string `gorm:"column:color;size:45"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	CreatedBy int64 `gorm:"column:created_by;omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;omitempty"`
	DeletedAt gorm.DeletedAt `swaggertype:"string" gorm:"index;column:deleted_at"`
	Clients	[]*ClientClientTag `gorm:"references:id;foreignKey:tag_id"`
}

func (*ClientTag) TableName() string {
	return "client_tags"
}

func (r *ClientTag) GetID() int64 {
	if r == nil {
		return 0
	}

	return r.ID
}

func (r *ClientTag) GetUUID() string {
	if r == nil {
		return ""
	}

	return r.UUID
}

func (r *ClientTag) GetName() string {
	if r == nil {
		return ""
	}

	return r.Name
}

func (r *ClientTag) GetColor() string {
	if r == nil {
		return "#FFFFFF"
	}

	if r.Color == "" {
		return "#FFFFFF"
	}

	return r.Color
}

func (r *ClientTag) GetCreatedAt() time.Time {
	if r == nil {
		return time.Now()
	}

	return r.CreatedAt
}

type ClientClientTag struct {
	ClientID 	int64 `gorm:"column:client_id;primaryKey"`
	Client 		*Client `gorm:"references:client_id;foreignKey:id;omitempty"`
	TagID 		int64 `gorm:"column:tag_id;primaryKey"`
	Tag 		*ClientTag `gorm:"references:tag_id;foreignKey:id;omitempty"`
	Color		string `gorm:"column:color;size:45"`
	CreatedBy 	int64 `gorm:"column:created_by;omitempty"`
	Creator   	*User `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt 	time.Time `gorm:"column:created_at;omitempty"`
}

func (*ClientClientTag) TableName() string {
	return "client_client_tags"
}

func (c *ClientClientTag) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}

func (c *ClientClientTag) GetTag() *ClientTag {
	if c == nil {
		return nil
	}

	return c.Tag
}