package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContactMailActivity interface {
	Create(ctx context.Context, data *entities.ContactMailActivity) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactMailActivity, error)
}

type ContactMailActivity struct {
}

func (ContactMailActivity) Create(ctx context.Context, data *entities.ContactMailActivity) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (ContactMailActivity) List(ctx context.Context, filters map[string]any) ([]*entities.ContactMailActivity, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.ContactMailActivity{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Unscoped().Preload("Creator").Preload("Mail").Preload("Contact").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
