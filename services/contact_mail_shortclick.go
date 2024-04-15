package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactMailShortClick interface {
	Create(ctx context.Context, data *entities.ContactMailShortClick) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactMailShortClick, error)
}

type ContactMailShortClick struct {
	Model models.IContactMailShortClick
}

func NewContactMailShortClick() IContactMailShortClick {
	return &ContactMailShortClick{
		Model: models.ContactMailShortClick{},
	}
}

func (p *ContactMailShortClick) Create(ctx context.Context, data *entities.ContactMailShortClick) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ContactMailShortClick) List(ctx context.Context, contactID int64) ([]*entities.ContactMailShortClick, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
