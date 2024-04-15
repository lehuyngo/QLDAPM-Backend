package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientTagActivity interface {
	Create(ctx context.Context, data *entities.ClientTagActivity) (int64, error)
	List(ctx context.Context, clientID int64) ([]*entities.ClientTagActivity, error)
}

type ClientTagActivity struct {
	Model models.IClientTagActivity
}

func NewClientTagActivity() IClientTagActivity {
	return &ClientTagActivity{
		Model: models.ClientTagActivity{},
	}
}

func (p *ClientTagActivity) Create(ctx context.Context, data *entities.ClientTagActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ClientTagActivity) List(ctx context.Context, clientID int64) ([]*entities.ClientTagActivity, error) {
	return p.Model.List(ctx, map[string]any{"client_id": clientID})
}
