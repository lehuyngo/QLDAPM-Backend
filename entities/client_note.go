package entities

type ClientNote struct {
	ID       int64   `gorm:"column:id;primaryKey"`
	ClientID int64   `gorm:"column:client_id;omitempty"`
	Client   *Client `gorm:"references:client_id;foreignKey:id;omitempty"`
	Title    string  `gorm:"column:title;size:255"`
	Content  string  `gorm:"column:content;size:4096"`
	Color    string  `gorm:"column:color;size:45"`
	Status   Status  `gorm:"column:status"`
	Base
}

func (*ClientNote) TableName() string {
	return "client_notes"
}

func (r *ClientNote) GetID() int64 {
	if r == nil {
		return 0
	}

	return r.ID
}

func (c *ClientNote) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}

func (r *ClientNote) GetTitle() string {
	if r == nil {
		return ""
	}

	return r.Title
}

func (r *ClientNote) GetContent() string {
	if r == nil {
		return ""
	}

	return r.Content
}

func (r *ClientNote) GetColor() string {
	if r == nil {
		return ""
	}

	return r.Color
}
