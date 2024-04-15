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

type IMeeting interface {
	Create(ctx context.Context, data *entities.Meeting) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Meeting, error)
	UpdateLastActiveTime(ctx context.Context, id int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.Meeting, error)
}

type Meeting struct {
}

func (Meeting) Create(ctx context.Context, data *entities.Meeting) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}

func (Meeting) ReadByUUID(ctx context.Context, uuid string) (*entities.Meeting, error) {
	result := &entities.Meeting{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	err := db.Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (Meeting) UpdateLastActiveTime(ctx context.Context, id int64) error {
	return clients.MySQLClient.WithContext(ctx).Model(&entities.Meeting{}).Where("id = ?", id).Update("last_active_time", time.Now().UnixMilli()).Error
}

func (Meeting) List(ctx context.Context, filters map[string]any) ([]*entities.Meeting, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.Meeting{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").
		Preload("Attendees.User").
		Preload("Attendees.Contact").
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
