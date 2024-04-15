package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IContactProject interface {
	Create(ctx context.Context, data *entities.ContactProject) error
	CreateBatch(ctx context.Context, data []entities.ContactProject) error
	Delete(ctx context.Context, contactID, projectID int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactProject, error)
}

type ContactProject struct {
}

func (ContactProject) Create(ctx context.Context, data *entities.ContactProject) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(data).Error
}

func (ContactProject) CreateBatch(ctx context.Context, data []entities.ContactProject) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(&data).Error
}

func (ContactProject) Delete(ctx context.Context, contactID, projectID int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("contact_id = ? AND project_id = ?", contactID, projectID).Delete(&entities.ContactProject{}).Error
}

func (ContactProject) List(ctx context.Context, filters map[string]any) ([]*entities.ContactProject, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ContactProject{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Contact").Preload("Project").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
