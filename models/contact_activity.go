package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContactActivity interface {
	Create(ctx context.Context, data *entities.ContactActivity) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactActivity, error)
}

type ContactActivity struct {
}

func (ContactActivity) Create(ctx context.Context, data *entities.ContactActivity) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (ContactActivity) List(ctx context.Context, filters map[string]any) ([]*entities.ContactActivity, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.ContactActivity{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").Preload("Contact").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
