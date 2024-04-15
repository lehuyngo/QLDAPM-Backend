package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IBatchMail interface {
	Create(ctx context.Context, data *entities.BatchMail) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.BatchMail, error)
	UpdateStatus(ctx context.Context, id int64, status entities.MailStatus) error
	List(ctx context.Context, filters map[string]any) ([]*entities.BatchMail, error)
}

type BatchMail struct {
}

func (BatchMail) Create(ctx context.Context, data *entities.BatchMail) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.ID, err
}

func (BatchMail) ReadByUUID(ctx context.Context, uuid string) (*entities.BatchMail, error) {
	result := &entities.BatchMail{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.BatchMail{})
	err := db.Preload("Creator").
		Preload("Receivers.Contact").
		Preload("CarbonCopies.User").
		Preload("AttachFiles").
		Preload("URLs").
		Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (BatchMail) UpdateStatus(ctx context.Context, id int64, status entities.MailStatus) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.BatchMail{}).Where("id = ?", id).Update("status", int(status)).Error
}

func (BatchMail) List(ctx context.Context, filters map[string]any) ([]*entities.BatchMail, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.BatchMail{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Receivers.Contact").
		Preload("CarbonCopies.User").
		Preload("AttachFiles").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
