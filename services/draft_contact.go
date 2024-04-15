package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IDraftContact interface {
	Create(ctx context.Context, data *entities.DraftContact) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.DraftContact, error)
	UpdateContact(ctx context.Context, id int64, contactID int64) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, orgID int64) ([]*entities.DraftContact, error)
	Restore(ctx context.Context, orgID int64) error
}

type DraftContact struct {
	Model models.IDraftContact
}

func NewDraftContact() IDraftContact{
	return &DraftContact{
		Model: models.DraftContact{},
	}
}

func (p *DraftContact) Create(ctx context.Context, data *entities.DraftContact) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *DraftContact) ReadByUUID(ctx context.Context, uuid string) (*entities.DraftContact, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *DraftContact) UpdateContact(ctx context.Context, id int64, contactID int64) error {
	return p.Model.UpdateField(ctx, id, "contact_id", contactID)
}

func (p *DraftContact) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *DraftContact) List(ctx context.Context, orgID int64) ([]*entities.DraftContact, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *DraftContact) Restore(ctx context.Context, orgID int64) error {
	return p.Model.Restore(ctx, orgID)
}