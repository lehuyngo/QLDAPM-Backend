package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IClientTag interface {
	Create(ctx context.Context, data *entities.ClientTag) (int64, error)
	First(ctx context.Context, filters map[string]any) (*entities.ClientTag, error)
	Update(ctx context.Context, data *entities.ClientTag) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ClientTag, error)
	Delete(ctx context.Context, id int64) error
}

type ClientTag struct {
}

func (ClientTag) Create(ctx context.Context, data *entities.ClientTag) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (ClientTag) First(ctx context.Context, filters map[string]any) (*entities.ClientTag, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := &entities.ClientTag{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.ClientTag{})
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Clients.Client").First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ClientTag) Update(ctx context.Context, data *entities.ClientTag) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit("Clients").Save(data).Where("id = ?", data.ID).Error
	})
}

func (ClientTag) List(ctx context.Context, filters map[string]any) ([]*entities.ClientTag, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.ClientTag{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Clients.Client").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ClientTag) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Select("Clients").Delete(&entities.ClientTag{ID: id}).Error
}
