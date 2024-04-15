package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IOrganization interface {
	Create(ctx context.Context, data *entities.Organization) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Organization, error)
	First(ctx context.Context) (*entities.Organization, error)
	Update(ctx context.Context, data *entities.Organization) error
}

type Organization struct {
}

func (Organization) Create(ctx context.Context, data *entities.Organization) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (Organization) Read(ctx context.Context, id int64) (*entities.Organization, error) {
	result := &entities.Organization{}
	db := clients.MySQLClient.WithContext(ctx)
	err := db.Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Organization) First(ctx context.Context) (*entities.Organization, error) {
	result := &entities.Organization{}
	err := clients.MySQLClient.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Organization) Update(ctx context.Context, data *entities.Organization) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Session(&gorm.Session{FullSaveAssociations: true}).Save(data).Where("id = ?", data.ID).Error
	})
}
