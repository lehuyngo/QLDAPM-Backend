package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IMail interface {
	Create(ctx context.Context, data *entities.Mail) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Mail, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.Mail, error)
}

type Mail struct {
}

func (Mail) Create(ctx context.Context, data *entities.Mail) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.ID, err
}

func (Mail) ReadByUUID(ctx context.Context, uuid string) (*entities.Mail, error) {
	result := &entities.Mail{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.Mail{})
	err := db.Preload("Creator").
		Preload("Receivers.Contact").
		Preload("Receivers.User").
		Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Mail) List(ctx context.Context, filters map[string]any) ([]*entities.Mail, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.Mail{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Receivers.Contact").
		Preload("Receivers.User").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
