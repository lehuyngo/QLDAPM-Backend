package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITrackedURL interface {
	Create(ctx context.Context, data *entities.TrackedURL) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.TrackedURL, error)
	First(ctx context.Context, filters map[string]any) (*entities.TrackedURL, error)
	UpdateStatus(ctx context.Context, id int64, status entities.ReadStatus) error
	Update(ctx context.Context, data *entities.TrackedURL) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.TrackedURL, error)
}

type TrackedURL struct {
}

func (TrackedURL) Create(ctx context.Context, data *entities.TrackedURL) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.ID, err
}

func (TrackedURL) ReadByUUID(ctx context.Context, uuid string) (*entities.TrackedURL, error) {
	result := &entities.TrackedURL{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.TrackedURL{})
	err := db.Preload("Creator").
		Preload("Mail").
		Preload("BatchMail").
		Preload("Contact").
		Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}


func (TrackedURL) First(ctx context.Context, filters map[string]any) (*entities.TrackedURL, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := &entities.TrackedURL{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Mail").
		Preload("BatchMail").
		Preload("Contact").
		WithContext(ctx).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (TrackedURL) UpdateStatus(ctx context.Context, id int64, status entities.ReadStatus) error {
	return clients.MySQLClient.WithContext(ctx).Omit(clause.Associations).Model(&entities.TrackedURL{}).Where("id = ?", id).Update("status", status.Value()).Error
}

func (TrackedURL) Update(ctx context.Context, data *entities.TrackedURL) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit(clause.Associations).Save(data).Where("id = ?", data.ID).Error
	})
}

func (TrackedURL) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.TrackedURL{}, id).Error
}

func (TrackedURL) List(ctx context.Context, filters map[string]any) ([]*entities.TrackedURL, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.TrackedURL{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Mail").
		Preload("BatchMail").
		Preload("Contact").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
