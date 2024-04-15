package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientAttachFile interface {
	Create(ctx context.Context, data *entities.ClientAttachFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.ClientAttachFile, error)
	List(ctx context.Context, clientID int64) ([]*entities.ClientAttachFile, error)
	Delete(ctx context.Context, clientID, id int64) error
}

type ClientAttachFile struct {
	Model models.IClientAttachFile
	ClientModel models.IClient
}

func NewClientAttachFile() IClientAttachFile {
	return &ClientAttachFile{
		Model: models.ClientAttachFile{},
		ClientModel: models.Client{},
	}
}

func (p *ClientAttachFile) Create(ctx context.Context, data *entities.ClientAttachFile) (string, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return "", err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)
	return result, nil
}

func (p *ClientAttachFile) Read(ctx context.Context, uuid string) (*entities.ClientAttachFile, error) {
	return p.Model.Read(ctx, uuid)
}

func (p *ClientAttachFile) List(ctx context.Context, clientID int64) ([]*entities.ClientAttachFile, error) {
	filters := map[string]any{"client_id": clientID}
	return p.Model.List(ctx, filters)
}

func (p *ClientAttachFile) Delete(ctx context.Context, clientID, id int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)
	return nil
}
