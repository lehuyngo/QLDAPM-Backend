package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IEventClickURL interface {
	Create(ctx context.Context, data *entities.EventClickURL) error
	List(ctx context.Context, orgID, minTime, maxTime int64, offset, limit int) ([]*entities.EventClickURL, error)
}

type EventClickURL struct {
}

func (EventClickURL) Create(ctx context.Context, data *entities.EventClickURL) error {
	return clients.MySQLClient.WithContext(ctx).Create(data).Error
}

func (EventClickURL) List(ctx context.Context, orgID, minTime, maxTime int64, offset, limit int) ([]*entities.EventClickURL, error) {
	result := []*entities.EventClickURL{}
	err := clients.MySQLClient.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Where("organization_id = ? AND read_time >= ? AND read_time <= ?", orgID, minTime, maxTime).Limit(limit).Offset(offset).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
