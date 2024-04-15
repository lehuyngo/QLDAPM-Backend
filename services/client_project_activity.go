package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientProjectActivity interface {
	Create(ctx context.Context, data *entities.ClientProjectActivity) (int64, error)
	List(ctx context.Context, clientID int64) ([]*entities.ClientProjectActivity, error)
}

type ClientProjectActivity struct {
	Model models.IClientProjectActivity
}

func NewClientProjectActivity() IClientProjectActivity {
	return &ClientProjectActivity{
		Model: models.ClientProjectActivity{},
	}
}

func (p *ClientProjectActivity) Create(ctx context.Context, data *entities.ClientProjectActivity) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *ClientProjectActivity) List(ctx context.Context, clientID int64) ([]*entities.ClientProjectActivity, error) {
	return p.Model.List(ctx, map[string]any{"client_id": clientID})
}
