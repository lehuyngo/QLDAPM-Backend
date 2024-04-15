package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContactTag interface {
	Create(ctx context.Context, data *entities.ContactTag) (int64, error)
	First(ctx context.Context, filters map[string]any) (*entities.ContactTag, error)
	Update(ctx context.Context, data *entities.ContactTag) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ContactTag, error)
	Delete(ctx context.Context, id int64) error
}

type ContactTag struct {
}

func (ContactTag) Create(ctx context.Context, data *entities.ContactTag) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	
	return data.ID, err
}
func (ContactTag) First(ctx context.Context, filters map[string]any) (*entities.ContactTag, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := &entities.ContactTag{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.ContactTag{})
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Contacts.Contact").First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactTag) Update(ctx context.Context, data *entities.ContactTag) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit("Contacts").Save(data).Where("id = ?", data.ID).Error
	})
}

func (ContactTag) List(ctx context.Context, filters map[string]any) ([]*entities.ContactTag, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ContactTag{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Contacts.Contact").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactTag) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("id = ?", id).Delete(&entities.ContactTag{}).Error
}