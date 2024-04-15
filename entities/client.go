package entities

type Client struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	FullName 		string `gorm:"column:fullname;size:200"`
	ShortName 		string `gorm:"column:shortname;size:200"`
	Code 			string `gorm:"column:code;size:200"`
	Fax 			string `gorm:"column:fax;size:200"`
	Website 		string `gorm:"column:website;size:500"`
	Phone			string `gorm:"column:phone;size:200"`
	Email			string `gorm:"column:email;size:200"`
	CompanySize		string `gorm:"column:company_size;size:200"`
	Address			string `gorm:"column:address;size:500"`
	Status   		Status `gorm:"column:status"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	LogoID			string `gorm:"column:logo_id;omitempty"`
	Logo   			*File `gorm:"references:logo_id;foreignKey:uuid;omitempty"`
	Notes 			[]*ClientNote `gorm:"references:id;foreignKey:client_id"`
	Tags 			[]*ClientClientTag `gorm:"references:id;foreignKey:client_id"`
	Projects 		[]*Project `gorm:"references:id;foreignKey:client_id"`
	Contacts 		[]*ClientContact `gorm:"references:id;foreignKey:client_id"`
	LastActiveTime	int64 `gorm:"column:last_active_time;omitempty"`
	Base
}

func (c *Client) Available() bool {
	if c == nil {
		return false
	}

	return c.Status > Inactive
}

func (*Client) TableName() string {
	return "clients"
}

func (c *Client) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *Client) GetFullName() string {
	if c == nil {
		return ""
	}

	return c.FullName
}

func (c *Client) GetShortName() string {
	if c == nil {
		return ""
	}

	return c.ShortName
}

func (c *Client) GetCode() string {
	if c == nil {
		return ""
	}

	return c.Code
}

func (c *Client) GetFax() string {
	if c == nil {
		return ""
	}

	return c.Fax
}

func (c *Client) GetEmail() string {
	if c == nil {
		return ""
	}

	return c.Email
}

func (c *Client) GetWebsite() string {
	if c == nil {
		return ""
	}

	return c.Website
}

func (c *Client) GetAddress() string {
	if c == nil {
		return ""
	}

	return c.Address
}

func (c *Client) GetPhone() string {
	if c == nil {
		return ""
	}

	return c.Phone
}

func (r *Client) GetLogo() *File {
	if r == nil {
		return nil
	}

	return r.Logo
}

func (r *Client) GetLastActiveTime() int64 {
	if r == nil {
		return 0
	}

	return r.LastActiveTime
}