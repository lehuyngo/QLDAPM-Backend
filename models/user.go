package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IUser interface {
	Create(ctx context.Context, data *entities.User) (int64, error)
	Read(ctx context.Context, id int64) (*entities.User, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.User, error)
	Update(ctx context.Context, data *entities.User) error
	Delete(ctx context.Context, id int64) error
	ReadByCondition(ctx context.Context, field string, value any) (*entities.User, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.User, error)
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.User, error)
}

type User struct {
}

func (User) Create(ctx context.Context, data *entities.User) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (User) Read(ctx context.Context, id int64) (*entities.User, error) {
	result := &entities.User{}
	db := clients.MySQLClient.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Table(result.TableName())
	err := db.Preload("Organization").Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (User) ReadByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	result := &entities.User{}
	db := clients.MySQLClient.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Table(result.TableName())
	err := db.Preload("Organization").Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (User) Update(ctx context.Context, data *entities.User) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Session(&gorm.Session{FullSaveAssociations: true}).Save(data).Where("id = ?", data.ID).Error
	})
}

func (User) Delete(ctx context.Context, id int64) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(&entities.User{}).Session(&gorm.Session{FullSaveAssociations: true}).Delete(&entities.User{}).Where("id = ?", id).Error
	})
}

func (User) ReadByCondition(ctx context.Context, field string, value any) (*entities.User, error) {
	result := &entities.User{}
	db := clients.MySQLClient.WithContext(ctx)
	err := db.Where(fmt.Sprintf("%s = ?", field), value).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (User) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.User, error) {
	result := []*entities.User{}
	err := clients.MySQLClient.WithContext(ctx).Where("uuid IN ?", uuids).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (User) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.User, error) {
	result := []*entities.User{}
	err := clients.MySQLClient.WithContext(ctx).Where("organization_id = ?", orgID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
