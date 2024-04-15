package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IStaticFile interface {
	Create(ctx context.Context, data *entities.StaticFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.StaticFile, error)
}

type StaticFile struct {
}

func (StaticFile) Create(ctx context.Context, data *entities.StaticFile) (string, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.UUID, err
}

func (StaticFile) Read(ctx context.Context, uuid string) (*entities.StaticFile, error) {
	result := &entities.StaticFile{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
