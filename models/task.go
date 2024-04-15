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

type ITask interface {
	Create(ctx context.Context, data *entities.Task) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.Task, error)
	Read(ctx context.Context, id int64) (*entities.Task, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Task, error)
	Update(ctx context.Context, data *entities.Task) error
	Delete(ctx context.Context, id int64) error
}

type Task struct {
}

func (Task) Create(ctx context.Context, data *entities.Task) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (Task) List(ctx context.Context, filters map[string]any) ([]*entities.Task, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.Task{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Project").Preload("Assignees.Assignee").Preload("AttachFiles.Creator").Preload("Creator").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Task) Read(ctx context.Context, id int64) (*entities.Task, error) {
	result := &entities.Task{}
	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Preload("Project").Preload("Assignees.Assignee").Preload("AttachFiles.Creator").Preload("Creator").Where("id = ?", id).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Task) ReadByUUID(ctx context.Context, uuid string) (*entities.Task, error) {
	result := &entities.Task{}
	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.Task{})
	err := db.Preload("Project").Preload("Assignees.Assignee").Preload("AttachFiles.Creator").Preload("Creator").Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Task) Update(ctx context.Context, data *entities.Task) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit(clause.Associations).Save(data).Where("id = ?", data.ID).Error
	})
}

func (Task) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.Task{}, id).Error
}
