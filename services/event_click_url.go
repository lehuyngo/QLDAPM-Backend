package services

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IEventClickURL interface {
	Create(ctx context.Context, data *entities.EventClickURL) error
	List(ctx context.Context, projectID, minTime, maxTime int64, offset, limit int) ([]*entities.EventClickURL, error)
}

type EventClickURL struct {
	Model models.IEventClickURL
}

func NewEventClickURL() IEventClickURL {
	return &EventClickURL{
		Model: models.EventClickURL{},
	}
}

func (p *EventClickURL) Create(ctx context.Context, data *entities.EventClickURL) error {
	return p.Model.Create(ctx, data)
}

func (p *EventClickURL) List(ctx context.Context, projectID, minTime, maxTime int64, offset, limit int) ([]*entities.EventClickURL, error) {
	return p.Model.List(ctx, projectID, minTime, maxTime, offset, limit)
}
