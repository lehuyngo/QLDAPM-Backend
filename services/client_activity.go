package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientActivity interface {
	Create(ctx context.Context, data *entities.ClientActivity) (int64, error)
	List(ctx context.Context, clientID int64) ([]*entities.ClientActivity, error)
}

type ClientActivity struct {
	Model models.IClientActivity
}

func NewClientActivity() IClientActivity {
	return &ClientActivity{
		Model: models.ClientActivity{},
	}
}

func (p *ClientActivity) Create(ctx context.Context, data *entities.ClientActivity) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)

}

func (p *ClientActivity) List(ctx context.Context, clientID int64) ([]*entities.ClientActivity, error) {
	return p.Model.List(ctx, map[string]any{"client_id": clientID})
}
