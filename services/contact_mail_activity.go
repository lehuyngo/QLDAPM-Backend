package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactMailActivity interface {
	Create(ctx context.Context, data *entities.ContactMailActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactMailActivity, error)
}

type ContactMailActivity struct {
	Model models.IContactMailActivity
}

func NewContactMailActivity() IContactMailActivity {
	return &ContactMailActivity{
		Model: models.ContactMailActivity{},
	}
}

func (p *ContactMailActivity) Create(ctx context.Context, data *entities.ContactMailActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ContactMailActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactMailActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
