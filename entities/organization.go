package entities

type Organization struct {
	ID       		int64 `gorm:"column:id;primaryKey"`
	DisplayName 	string `gorm:"column:displayname;size:200"`
	Base
}

func (*Organization) TableName() string {
	return "organizations"
}

func (o *Organization) GetID() int64 {
	if o == nil {
		return 0
	}

	return o.ID
}