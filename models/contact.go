package models

import (
	"context"
	"fmt"
	"time"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IContact interface {
	Create(ctx context.Context, data *entities.Contact) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Contact, error)
	First(ctx context.Context, filters map[string]any) (*entities.Contact, error)
	Update(ctx context.Context, data *entities.Contact) error
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	UpdateLastActiveTime(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter map[string]any) ([]*entities.Contact, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contact, error)
}

type Contact struct {
}

func (Contact) Create(ctx context.Context, data *entities.Contact) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (Contact) ReadByUUID(ctx context.Context, uuid string) (*entities.Contact, error) {
	result := &entities.Contact{}
	db := clients.MySQLClient.WithContext(ctx).Model(&entities.Contact{})
	err := db.Preload("Creator.Organization").
		Preload("Avatar").
		Preload("NameCard").
		Preload("Notes").
		Preload("Tags.Tag").
		Where("uuid = ?", uuid).
		First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Contact) First(ctx context.Context, filters map[string]any) (*entities.Contact, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := &entities.Contact{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Table(result.TableName())
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator.Organization").
		Preload("Avatar").
		Preload("NameCard").
		Preload("Notes").
		Preload("Tags.Tag").
		First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Contact) Update(ctx context.Context, data *entities.Contact) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Save(data).Where("id = ?", data.ID).Error
	})
}

func (Contact) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Contact{}).Where("id = ?", id).Update(field, value).Error
}

func (Contact) UpdateLastActiveTime(ctx context.Context, id int64) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Contact{}).Where("id = ?", id).Update("last_active_time", time.Now().UnixMilli()).Error
}

func (Contact) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.Contact{}, id).Error
}

func (Contact) List(ctx context.Context, filters map[string]any) ([]*entities.Contact, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.Contact{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Avatar").
		Preload("NameCard").
		Preload("Notes").
		Preload("Tags.Tag").
		Preload("Clients.Client", "status >= (?)", entities.Active).
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (Contact) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contact, error) {
	result := []*entities.Contact{}
	err := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Where("uuid IN ?", uuids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}