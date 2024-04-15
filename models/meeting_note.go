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

type IMeetingNote interface {
	Create(ctx context.Context, data *entities.MeetingNote) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingNote, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.MeetingNote, error)
	Update(ctx context.Context, data *entities.MeetingNote) error
	Delete(ctx context.Context, id int64) error
}

type MeetingNote struct {
}

func (MeetingNote) Create(ctx context.Context, data *entities.MeetingNote) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (MeetingNote) ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingNote, error) {
	result := &entities.MeetingNote{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Preload("Project").Preload("Contributors.User").Preload("Contributors.Contact").Preload("Editors.Editor").Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingNote) List(ctx context.Context, filters map[string]any) ([]*entities.MeetingNote, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.MeetingNote{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Contributors.User").
		Preload("Contributors.Contact").
		Preload("Editors.Editor").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingNote) Update(ctx context.Context, data *entities.MeetingNote) error {
	return clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(data).Omit(clause.Associations).Save(data).Where("id = ?", data.ID).Error
	})
}

func (MeetingNote) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.MeetingNote{}, id).Error
}
