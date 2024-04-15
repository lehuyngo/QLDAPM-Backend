package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactActivity interface {
	Create(ctx context.Context, data *entities.ContactActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactActivity, error)
}

type ContactActivity struct {
	Model models.IContactActivity
}

func NewContactActivity() IContactActivity {
	return &ContactActivity{
		Model: models.ContactActivity{},
	}
}

func (p *ContactActivity) Create(ctx context.Context, data *entities.ContactActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ContactActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
