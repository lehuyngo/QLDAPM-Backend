package entities

type Contact struct {
	ID   			int64 `gorm:"column:id;primaryKey"`
	FullName 		string `gorm:"column:fullname;size:200"`
	ShortName 		string `gorm:"column:shortname;size:200"`
	JobTitle 		string `gorm:"column:job_title;size:200"`
	Phone			string `gorm:"column:phone;size:200"`
	Email			string `gorm:"column:email;size:200"`
	Gender   		Gender `gorm:"column:gender"`
	Status   		Status `gorm:"column:status"`
	OrganizationID	int64 `gorm:"column:organization_id;omitempty"`
	AvatarID		string `gorm:"column:avatar_id;omitempty"`
	Avatar   		*File `gorm:"references:avatar_id;foreignKey:uuid;omitempty"`
	NameCardID		string `gorm:"column:name_card_id;omitempty"`
	NameCard   		*File `gorm:"references:name_card_id;foreignKey:uuid;omitempty"`
	Notes 			[]*ContactNote `gorm:"references:id;foreignKey:contact_id"`
	Tags 			[]*ContactContactTag `gorm:"references:id;foreignKey:contact_id"`
	Clients 		[]*ClientContact `gorm:"references:id;foreignKey:contact_id"`
	LastActiveTime	int64 `gorm:"column:last_active_time;omitempty"`
	Birthday 		int64 `gorm:"column:birthday;omitempty"`
	*Base
}

func (c *Contact) Available() bool {
	if c == nil {
		return false
	}

	return c.Status > Inactive
}

func (*Contact) TableName() string {
	return "contacts"
}

func (c *Contact) GetID() int64 {
	if c == nil {
		return 0
	}

	return c.ID
}

func (c *Contact) GetFullName() string {
	if c == nil {
		return ""
	}

	return c.FullName
}

func (c *Contact) GetShortName() string {
	if c == nil {
		return ""
	}

	return c.ShortName
}

func (c *Contact) GetPhone() string {
	if c == nil {
		return ""
	}

	return c.Phone
}

func (c *Contact) GetEmail() string {
	if c == nil {
		return ""
	}

	return c.Email
}

func (c *Contact) GetJobTitle() string {
	if c == nil {
		return ""
	}

	return c.JobTitle
}

func (c *Contact) GetBitrhDay() int64 {
	if c == nil {
		return 0
	}

	return c.Birthday
}

func (c *Contact) GetGender() Gender {
	if c == nil {
		return Male
	}

	return c.Gender
}

func (r *Contact) GetAvatar() *File {
	if r == nil {
		return nil
	}

	return r.Avatar
}

func (r *Contact) GetNameCard() *File {
	if r == nil {
		return nil
	}

	return r.NameCard
}