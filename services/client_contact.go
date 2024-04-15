package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientContact interface {
	Create(ctx context.Context, clientID int64, data *entities.Contact) error
	Add(ctx context.Context, data *entities.ClientContact) error
	AddBatch(ctx context.Context, data []entities.ClientContact) error
	Delete(ctx context.Context, userID, clientID, contactID int64) error
	List(ctx context.Context, clientID int64) ([]*entities.ClientContact, error)
}

type ClientContact struct {
	Model models.IClientContact
	ClientModel models.IClient
	ContactModel models.IContact
	ContactClientActivity models.IContactClientActivity // đều dùng chung model này 
}

func NewClientContact() IClientContact {
	return &ClientContact{
		Model: models.ClientContact{},
		ClientModel: models.Client{},
		ContactModel: models.Contact{},
		ContactClientActivity: models.ContactClientActivity{},
	}
}

func (p *ClientContact) Create(ctx context.Context, clientID int64, data *entities.Contact) error {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	contactID, err := p.ContactModel.Create(ctx, data)
	if err != nil {
		return err
	}

	err = p.Model.Create(ctx, &entities.ClientContact{
		ClientID: clientID,
		ContactID: contactID,
		CreatedBy: data.CreatedBy,
	})
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)
	p.ContactClientActivity.Create(ctx, &entities.ContactClientActivity{
		ClientID: clientID,
		ContactID: contactID,
		Type: entities.ActivityCreated,
		CreatedBy: data.CreatedBy,
	})
	return nil
}

func (p *ClientContact) Add(ctx context.Context, data *entities.ClientContact) error {
	err := p.Model.Create(ctx, data)
	if err != nil {
		return err
	}
	p.ContactClientActivity.Create(ctx, &entities.ContactClientActivity{
		ClientID: data.ClientID,
		ContactID: data.ContactID,
		Type: entities.ActivityCreated,
		CreatedBy: data.CreatedBy,
	})
	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)
	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)
	return nil
}

func (p *ClientContact) AddBatch(ctx context.Context, data []entities.ClientContact) error {
	err := p.Model.CreateBatch(ctx, data)
	if err != nil {
		return err
	}

	for _, val := range data {
		p.ContactClientActivity.Create(ctx, &entities.ContactClientActivity{
			ClientID: val.ClientID,
			ContactID: val.ContactID,
			Type: entities.ActivityCreated,
			CreatedBy: val.CreatedBy,
		})
		p.ClientModel.UpdateLastActiveTime(ctx, val.ClientID)
		p.ContactModel.UpdateLastActiveTime(ctx, val.ContactID)
	}
	return nil
}

func (p *ClientContact) Delete(ctx context.Context, userID,clientID, contactID int64) error {
	err := p.Model.Delete(ctx, clientID, contactID)
	if err != nil {
		return err
	}
	p.ContactClientActivity.Create(ctx, &entities.ContactClientActivity{
		ClientID: clientID,
		ContactID: contactID,
		Type: entities.ActivityDeleted,
		CreatedBy: userID,
	})

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)
	p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	return nil
}

func (p *ClientContact) List(ctx context.Context, clientID int64) ([]*entities.ClientContact, error) {
	filters := map[string]any{"client_id": clientID}
	return  p.Model.List(ctx, filters)
}
