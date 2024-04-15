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

type IClient interface {
	Create(ctx context.Context, data *entities.Client) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Client, error)
	First(ctx context.Context, filters map[string]any) (*entities.Client, error)
	Update(ctx context.Context, data *entities.Client) error
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	UpdateLastActiveTime(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter map[string]any) ([]*entities.Client, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Client, error)
}

type Client struct {
}

func (Client) Create(ctx context.Context, data *entities.Client) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (Client) ReadByUUID(ctx context.Context, uuid string) (*entities.Client, error) {
	result := &entities.Client{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.Client{})
	err := db.Preload("Creator.Organization").
		Preload("Logo").
		Preload("Notes").
		Preload("Tags.Tag").
		Preload("Projects", "status >= (?)", entities.Active).
		Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Client) First(ctx context.Context, filters map[string]any) (*entities.Client, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := &entities.Client{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Table(result.TableName())
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator.Organization").
		Preload("Logo").
		Preload("Notes").
		Preload("Tags.Tag").
		Preload("Projects", "status >= (?)", entities.Active).
		WithContext(ctx).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Client) Update(ctx context.Context, data *entities.Client) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Save(data).Where("id = ?", data.ID).Error
	})
}

func (Client) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Client{}).Where("id = ?", id).Update(field, value).Error
}

func (Client) UpdateLastActiveTime(ctx context.Context, id int64) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Client{}).Where("id = ?", id).Update("last_active_time", time.Now().UnixMilli()).Error
}

func (Client) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.Client{}, id).Error
}

func (Client) List(ctx context.Context, filters map[string]any) ([]*entities.Client, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.Client{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Logo").
		Preload("Notes").
		Preload("Tags.Tag").
		Preload("Projects", "status >= (?)", entities.Active).
		Preload("Contacts.Contact").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (Client) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Client, error) {
	result := []*entities.Client{}
	err := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Where("uuid IN ?", uuids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}