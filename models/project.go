package models

import (
	"context"
	"fmt"
	"time"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IProject interface {
	Create(ctx context.Context, data *entities.Project) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Project, error)
	First(ctx context.Context, filters map[string]any) (*entities.Project, error)
	Update(ctx context.Context, data *entities.Project) error
	UpdateField(ctx context.Context, id int64, field string, value interface{}) error
	UpdateLastActiveTime(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter map[string]any) ([]*entities.Project, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Project, error)
}

type Project struct {
}

func (Project) Create(ctx context.Context, data *entities.Project) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (Project) ReadByUUID(ctx context.Context, uuid string) (*entities.Project, error) {
	result := &entities.Project{}
	db := clients.MySQLClient.WithContext(ctx).Model(&entities.Project{})
	err := db.Preload("Client").Preload("Contacts.Contact").
		Where("uuid = ?", uuid).
		First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Project) First(ctx context.Context, filters map[string]any) (*entities.Project, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := &entities.Project{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Table(result.TableName())
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Client").Preload("Contacts.Contact").First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Project) Update(ctx context.Context, data *entities.Project) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit(clause.Associations).Save(data).Where("id = ?", data.ID).Error
	})
}

func (Project) UpdateField(ctx context.Context, id int64, field string, value interface{}) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Project{}).Where("id = ?", id).Update(field, value).Error
}

func (Project) UpdateLastActiveTime(ctx context.Context, id int64) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Project{}).Where("id = ?", id).Update("last_active_time", time.Now().UnixMilli()).Error
}

func (Project) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.Project{}, id).Error
}

func (Project) List(ctx context.Context, filters map[string]any) ([]*entities.Project, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.Project{}

	db := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Client").Preload("Contacts.Contact").Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

func (Project) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Project, error) {
	result := []*entities.Project{}
	err := clients.MySQLClient.WithContext(ctx).Scopes(NotInTrash).Where("uuid IN ?", uuids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	
	return result, nil
}
