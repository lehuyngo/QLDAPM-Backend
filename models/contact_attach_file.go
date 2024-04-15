package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContactAttachFile interface {
	Create(ctx context.Context, data *entities.ContactAttachFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.ContactAttachFile, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactAttachFile, error)
	Delete(ctx context.Context, id int64) error
}

type ContactAttachFile struct {
}

func (ContactAttachFile) Create(ctx context.Context, data *entities.ContactAttachFile) (string, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.UUID, err
}

func (ContactAttachFile) Read(ctx context.Context, uuid string) (*entities.ContactAttachFile, error) {
	result := &entities.ContactAttachFile{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactAttachFile) List(ctx context.Context, filters map[string]any) ([]*entities.ContactAttachFile, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ContactAttachFile{}
	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactAttachFile) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.ContactAttachFile{}).Error
}
