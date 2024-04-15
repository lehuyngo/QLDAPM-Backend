package services

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactClientActivity interface {
	Create(ctx context.Context, data *entities.ContactClientActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactClientActivity, error)
}

type ContactClientActivity struct {
	Model models.IContactClientActivity
}

func NewContactClientActivity() IContactClientActivity {
	return &ContactClientActivity{
		Model: models.ContactClientActivity{},
	}
}

func (p *ContactClientActivity) Create(ctx context.Context, data *entities.ContactClientActivity) (int64, error) {
	return p.Model.Create(ctx, data)
}

func (p *ContactClientActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactClientActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
