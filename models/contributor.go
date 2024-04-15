package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContributor interface {
	Create(ctx context.Context, data *entities.Contributor) (int64, error)
	CreateBatch(ctx context.Context, data []*entities.Contributor) error
	ReadByUUID(ctx context.Context, uuid string) (*entities.Contributor, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contributor, error)
	Delete(ctx context.Context, id int64) error
	DeleteBatch(ctx context.Context, uuids []string) error
}

type Contributor struct {
}

func (Contributor) Create(ctx context.Context, data *entities.Contributor) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (Contributor) CreateBatch(ctx context.Context, data []*entities.Contributor) error {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return err
}

func (Contributor) ReadByUUID(ctx context.Context, uuid string) (*entities.Contributor, error) {
	db := clients.MySQLClient.WithContext(ctx)
	var result entities.Contributor
	err := db.Where("uuid = ?", uuid).First(&result).Error
	return &result, err
}

func (Contributor) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contributor, error) {
	db := clients.MySQLClient.WithContext(ctx)
	var result []*entities.Contributor
	err := db.Where("uuid IN ?", uuids).Find(&result).Error
	return result, err
}

func (Contributor) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.Contributor{}).Error
}

func (Contributor) DeleteBatch(ctx context.Context, uuids []string) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("uuid IN ?", uuids).Delete(&entities.Contributor{}).Error
}
