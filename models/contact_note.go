package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IContactNote interface {
	Create(ctx context.Context, data *entities.ContactNote) (int64, error)
	Read(ctx context.Context, id int64) (*entities.ContactNote, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.ContactNote, error)
	Update(ctx context.Context, data *entities.ContactNote) error
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactNote, error)
}

type ContactNote struct {
}

func (ContactNote) Create(ctx context.Context, data *entities.ContactNote) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})
	
	return data.ID, err
}

func (ContactNote) Read(ctx context.Context, id int64) (*entities.ContactNote, error) {
	result := &entities.ContactNote{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactNote) ReadByUUID(ctx context.Context, uuid string) (*entities.ContactNote, error) {
	result := &entities.ContactNote{}
	db := clients.MySQLClient.WithContext(ctx)
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactNote) Update(ctx context.Context, data *entities.ContactNote) error {
	db := clients.MySQLClient.WithContext(ctx).Model(data)
	return db.Omit(clause.Associations).Save(data).Where("id = ?", data.ID).Error
}

func (ContactNote) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	db := clients.MySQLClient.WithContext(ctx).Model(&ContactNote{})
	return db.Where("id = ?", id).Update(field, value).Error
}

func (ContactNote) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.ContactNote{}).Error
}

func (ContactNote) List(ctx context.Context, filters map[string]any) ([]*entities.ContactNote, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ContactNote{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}