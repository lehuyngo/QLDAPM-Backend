package entities

type DraftContact struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	FullName 		string `gorm:"column:fullname;size:200"`
	Phone			string `gorm:"column:phone;size:200"`
	Email			string `gorm:"column:email;size:200"`
	ClientName 		string `gorm:"column:client_name;size:200"`
	ClientWebsite	string `gorm:"column:client_website;size:200"`
	ClientAddress	string `gorm:"column:client_address;size:500"`
	NameCardID		string `gorm:"column:name_card_id;omitempty"`
	NameCard   		*File `gorm:"references:name_card_id;foreignKey:uuid;omitempty"`
	CompanyLogoID	string `gorm:"column:company_logo_id;omitempty"`
	CompanyLogo  	*File `gorm:"references:company_logo_id;foreignKey:uuid;omitempty"`
	ContactID 		int64 `gorm:"column:contact_id;omitempty"`
	Contact   		*Contact `gorm:"references:contact_id;foreignKey:id;omitempty"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	Base
}

func (*DraftContact) TableName() string {
	return "draft_contacts"
}

func (c *DraftContact) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *DraftContact) GetFullName() string {
	if c == nil {
		return ""
	}

	return c.FullName
}

func (c *DraftContact) GetPhone() string {
	if c == nil {
		return ""
	}

	return c.Phone
}

func (c *DraftContact) GetEmail() string {
	if c == nil {
		return ""
	}

	return c.Email
}

func (c *DraftContact) GetClientName() string {
	if c == nil {
		return ""
	}

	return c.ClientName
}

func (c *DraftContact) GetClientWebsite() string {
	if c == nil {
		return ""
	}

	return c.ClientWebsite
}

func (c *DraftContact) GetClientAddress() string {
	if c == nil {
		return ""
	}

	return c.ClientAddress
}

func (r *DraftContact) GetNameCard() *File {
	if r == nil {
		return nil
	}

	return r.NameCard
}

func (r *DraftContact) GetCompanyLogo() *File {
	if r == nil {
		return nil
	}

	return r.CompanyLogo
}