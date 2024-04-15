package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IClientProjectActivity interface {
	Create(ctx context.Context, data *entities.ClientProjectActivity) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.ClientProjectActivity, error)
}

type ClientProjectActivity struct {
}

func (ClientProjectActivity) Create(ctx context.Context, data *entities.ClientProjectActivity) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	return data.ID, err
}

func (ClientProjectActivity) List(ctx context.Context, filters map[string]any) ([]*entities.ClientProjectActivity, error) {
	result := []*entities.ClientProjectActivity{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(field, val)
	}

	err := db.Preload("Creator").Preload("Project").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
