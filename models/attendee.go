package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IAttendee interface {
	Create(ctx context.Context, data *entities.Attendee) (int64, error)
}

type Attendee struct {
}

func (Attendee) Create(ctx context.Context, data *entities.Attendee) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}
