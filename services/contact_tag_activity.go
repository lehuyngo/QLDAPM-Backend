package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactTagActivity interface {
	Create(ctx context.Context, data *entities.ContactTagActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactTagActivity, error)
}

type ContactTagActivity struct {
	Model models.IContactTagActivity
}

func NewContactTagActivity() IContactTagActivity {
	return &ContactTagActivity{
		Model: models.ContactTagActivity{},
	}
}

func (p *ContactTagActivity) Create(ctx context.Context, data *entities.ContactTagActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ContactTagActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactTagActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
