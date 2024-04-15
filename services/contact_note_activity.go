package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactNoteActivity interface {
	Create(ctx context.Context, data *entities.ContactNoteActivity) (int64, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactNoteActivity, error)
}

type ContactNoteActivity struct {
	Model models.IContactNoteActivity
}

func NewContactNoteActivity() IContactNoteActivity {
	return &ContactNoteActivity{
		Model: models.ContactNoteActivity{},
	}
}

func (p *ContactNoteActivity) Create(ctx context.Context, data *entities.ContactNoteActivity) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *ContactNoteActivity) List(ctx context.Context, contactID int64) ([]*entities.ContactNoteActivity, error) {
	return p.Model.List(ctx, map[string]any{"contact_id": contactID})
}
