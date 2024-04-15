package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientNoteActivity interface {
	Create(ctx context.Context, data *entities.ClientNoteActivity) (int64, error)
	List(ctx context.Context, clientID int64) ([]*entities.ClientNoteActivity, error)
}

type ClientNoteActivity struct {
	Model models.IClientNoteActivity
}

func NewClientNoteActivity() IClientNoteActivity {
	return &ClientNoteActivity{
		Model: models.ClientNoteActivity{},
	}
}

func (p *ClientNoteActivity) Create(ctx context.Context, data *entities.ClientNoteActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ClientNoteActivity) List(ctx context.Context, clientID int64) ([]*entities.ClientNoteActivity, error) {
	return p.Model.List(ctx, map[string]any{"client_id": clientID})
}
