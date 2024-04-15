package entities

type Project struct {
	ID             int64             `gorm:"column:id;primaryKey"`
	FullName       string            `gorm:"column:fullname;size:200"`
	ShortName      string            `gorm:"column:shortname;size:200"`
	Code           string            `gorm:"column:code;size:200"`
	Status         Status            `gorm:"column:status"`
	ProjectStatus  ProjectStatus     `gorm:"column:project_status"`
	ClientID       int64             `gorm:"column:client_id;omitempty"`
	Client         *Client           `gorm:"references:client_id;foreignKey:id;omitempty"`
	OrganizationID int64             `gorm:"column:organization_id;omitempty"`
	LastActiveTime int64             `gorm:"column:last_active_time;omitempty"`
	Contacts       []*ContactProject `gorm:"references:id;foreignKey:project_id"`
	Base
}

func (c *Project) Available() bool {
	if c == nil {
		return false
	}

	return c.Status > Inactive
}

func (*Project) TableName() string {
	return "projects"
}

func (c *Project) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *Project) GetFullName() string {
	if c == nil {
		return ""
	}

	return c.FullName
}

func (c *Project) GetShortName() string {
	if c == nil {
		return ""
	}

	return c.ShortName
}

func (c *Project) GetCode() string {
	if c == nil {
		return ""
	}

	return c.Code
}

func (c *Project) GetClient() *Client {
	if c == nil {
		return nil
	}

	return c.Client
}

func (c *Project) GetLastActiveTime() int64 {
	if c == nil {
		return 0
	}

	return c.LastActiveTime
}

func (c *Project) GetOrganizationID() int64 {
	if c == nil {
		return 0
	}

	return c.OrganizationID
}
