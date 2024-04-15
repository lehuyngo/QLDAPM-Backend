package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactProjectActivity interface {
	Create(ctx context.Context, data *entities.ContactProjectActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactProjectActivity, error)
}

type ContactProjectActivity struct {
	Model models.IContactProjectActivity
}

func NewContactProjectActivity() IContactProjectActivity {
	return &ContactProjectActivity{
		Model: models.ContactProjectActivity{},
	}
}

func (p *ContactProjectActivity) Create(ctx context.Context, data *entities.ContactProjectActivity) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *ContactProjectActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactProjectActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
