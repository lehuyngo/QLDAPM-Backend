package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IClientNoteActivity interface {
	Create(ctx context.Context, data *entities.ClientNoteActivity) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.ClientNoteActivity, error)
}

type ClientNoteActivity struct {
}

func (ClientNoteActivity) Create(ctx context.Context, data *entities.ClientNoteActivity) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (ClientNoteActivity) List(ctx context.Context, filters map[string]any) ([]*entities.ClientNoteActivity, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.ClientNoteActivity{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Unscoped().Preload("Creator").Preload("Note").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
