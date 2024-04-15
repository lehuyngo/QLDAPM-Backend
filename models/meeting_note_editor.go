package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IMeetingNoteEditor interface {
	Create(ctx context.Context, data *entities.MeetingNoteEditor) (int64, error)
}

type MeetingNoteEditor struct {
}

func (MeetingNoteEditor) Create(ctx context.Context, data *entities.MeetingNoteEditor) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})

	return data.ID, err
}
