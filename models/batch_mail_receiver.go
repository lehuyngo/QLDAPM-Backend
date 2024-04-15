package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IBatchMailReceiver interface {
	Read(ctx context.Context, id int64) (*entities.BatchMailReceiver, error)
}

type BatchMailReceiver struct {
}

func (BatchMailReceiver) Read(ctx context.Context, id int64) (*entities.BatchMailReceiver, error) {
	result := &entities.BatchMailReceiver{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.BatchMailReceiver{})
	err := db.Preload("Mail.CarbonCopies").
		Preload("Mail.AttachFiles").
		Preload("Mail.URLs").
		Preload("Contact").
		Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (BatchMailReceiver) Update(ctx context.Context, data *entities.BatchMailReceiver) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Save(data).Where("id = ?", data.ID).Error
	})
}