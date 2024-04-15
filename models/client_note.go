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

type IClientNote interface {
	Create(ctx context.Context, data *entities.ClientNote) (int64, error)
	Read(ctx context.Context, id int64) (*entities.ClientNote, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.ClientNote, error)
	Update(ctx context.Context, data *entities.ClientNote) error
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ClientNote, error)
}

type ClientNote struct {
}

func (ClientNote) Create(ctx context.Context, data *entities.ClientNote) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})
	
	return data.ID, err
}

func (ClientNote) Read(ctx context.Context, id int64) (*entities.ClientNote, error) {
	result := &entities.ClientNote{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ClientNote) ReadByUUID(ctx context.Context, uuid string) (*entities.ClientNote, error) {
	result := &entities.ClientNote{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ClientNote) Update(ctx context.Context, data *entities.ClientNote) error {
	db := clients.MySQLClient.WithContext(ctx).Model(data)
	return db.Omit(clause.Associations).Where("id = ?", data.ID).Save(data).Error
}

func (ClientNote) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	db := clients.MySQLClient.WithContext(ctx).Model(&ClientNote{})
	return db.Where("id = ?", id).Update(field, value).Error
}

func (ClientNote) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.ClientNote{}).Error
}

func (ClientNote) List(ctx context.Context, filters map[string]any) ([]*entities.ClientNote, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ClientNote{}

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