package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type ITaskAttachFile interface {
	Create(ctx context.Context, data *entities.TaskAttachFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.TaskAttachFile, error)
	Delete(ctx context.Context, id int64) error
}

type TaskAttachFile struct {
}

func (TaskAttachFile) Create(ctx context.Context, data *entities.TaskAttachFile) (string, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.UUID, err
}

func (TaskAttachFile) Read(ctx context.Context, uuid string) (*entities.TaskAttachFile, error) {
	result := &entities.TaskAttachFile{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (TaskAttachFile) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.TaskAttachFile{}).Error
}
