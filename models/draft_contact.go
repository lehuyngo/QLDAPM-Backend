package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IDraftContact interface {
	Create(ctx context.Context, data *entities.DraftContact) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.DraftContact, error)
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.DraftContact, error)
	Restore(ctx context.Context, orgID int64) error
}

type DraftContact struct {
}

func (DraftContact) Create(ctx context.Context, data *entities.DraftContact) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (DraftContact) ReadByUUID(ctx context.Context, uuid string) (*entities.DraftContact, error) {
	result := &entities.DraftContact{}
	db := clients.MySQLClient.WithContext(ctx).Model(&entities.DraftContact{})
	err := db.Preload("NameCard").Preload("CompanyLogo").Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (DraftContact) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.DraftContact{}).Where("id = ?", id).Update(field, value).Error
}

func (DraftContact) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.DraftContact{}, id).Error
}

func (DraftContact) List(ctx context.Context, filters map[string]any) ([]*entities.DraftContact, error) {
	result := []*entities.DraftContact{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("NameCard").Preload("CompanyLogo").Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (DraftContact) Restore(ctx context.Context, orgID int64) error {
	return clients.MySQLClient.WithContext(ctx).Unscoped().Model(&entities.DraftContact{}).Where("organization_id = ?", orgID).Select("deleted_at").Updates(map[string]interface{}{"deleted_at": gorm.Expr("NULL")}).Error
}