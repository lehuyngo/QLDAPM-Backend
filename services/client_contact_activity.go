package services

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientContactActivity interface {
	Create(ctx context.Context, data *entities.ContactClientActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactClientActivity, error)
}

type ClientContactActivity struct {
	Model models.IContactClientActivity
}

func NewClientContactActivity() IClientContactActivity {
	return &ClientContactActivity{
		Model: models.ContactClientActivity{},
	}
}

func (p *ClientContactActivity) Create(ctx context.Context, data *entities.ContactClientActivity) (int64, error) {
	return p.Model.Create(ctx, data)
}

func (p *ClientContactActivity) List(ctx context.Context, clientID int64) ([]*entities.ContactClientActivity, error) {
	return p.Model.List(ctx, map[string]any{"client_id": clientID})
}
