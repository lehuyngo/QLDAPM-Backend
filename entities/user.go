package entities

type MigrateUser struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	DisplayName  string `gorm:"column:displayname;size:200"`
	Username     string `gorm:"column:username;size:200"`
	PasswordHash string `gorm:"column:password_hash;size:200"`
	Email        string `gorm:"column:email;size:200"`
	Status       int    `gorm:"column:status"`
}

func (*MigrateUser) TableName() string {
	return "users"
}

type User struct {
	ID             int64         `gorm:"column:id;primaryKey"`
	DisplayName    string        `gorm:"column:displayname;size:200"`
	OrganizationID int64         `gorm:"column:organization_id;omitempty"`
	Organization   *Organization `gorm:"references:organization_id;foreignKey:id;omitempty"`
	Username       string        `gorm:"column:username;size:200"`
	PasswordHash   string        `gorm:"column:password_hash;size:200"`
	Email          string        `gorm:"column:email;size:200"`
	Status         int32         `gorm:"column:status"`
	Token          string        `gorm:"column:token;size:200"`
	Base
}

func (*User) TableName() string {
	return "users"
}

func (r *User) GetID() int64 {
	if r == nil {
		return 0
	}

	return r.ID
}

func (r *User) GetPasswordHash() string {
	if r == nil {
		return ""
	}

	return r.PasswordHash
}

func (r *User) GetOrganization() *Organization {
	if r == nil {
		return nil
	}

	return r.Organization
}

func (r *User) GetDisplayName() string {
	if r == nil {
		return ""
	}

	return r.DisplayName
}

func (r *User) GetEmail() string {
	if r == nil {
		return ""
	}

	return r.Email
}
