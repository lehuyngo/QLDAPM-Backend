package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactAttachFile interface {
	Create(ctx context.Context, data *entities.ContactAttachFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.ContactAttachFile, error)
	List(ctx context.Context, contactID int64) ([]*entities.ContactAttachFile, error)
	Delete(ctx context.Context, contactID, id int64) error
}

type ContactAttachFile struct {
	Model models.IContactAttachFile
	ClientModel models.IClient
}

func NewContactAttachFile() IContactAttachFile {
	return &ContactAttachFile{
		Model: models.ContactAttachFile{},
		ClientModel: models.Client{},
	}
}

func (p *ContactAttachFile) Create(ctx context.Context, data *entities.ContactAttachFile) (string, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return "", err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, data.ContactID)
	return result, nil
}

func (p *ContactAttachFile) Read(ctx context.Context, uuid string) (*entities.ContactAttachFile, error) {
	return p.Model.Read(ctx, uuid)
}

func (p *ContactAttachFile) List(ctx context.Context, contactID int64) ([]*entities.ContactAttachFile, error) {
	filters := map[string]any{"contact_id": contactID}
	return p.Model.List(ctx, filters)
}

func (p *ContactAttachFile) Delete(ctx context.Context, contactID, id int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, contactID)
	return nil
}
