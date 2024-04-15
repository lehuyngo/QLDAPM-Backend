package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactClient interface {
	Create(ctx context.Context, contactID int64, data *entities.Client) error
	Add(ctx context.Context, data *entities.ClientContact) error
	AddBatch(ctx context.Context, data []entities.ClientContact) error
	Delete(ctx context.Context, userID ,contactID, clientID int64) error
	List(ctx context.Context, contactID int64) ([]*entities.ClientContact, error)
}

type ContactClient struct {
	Model         models.IClientContact
	ClientModel   models.IClient
	ContactModel  models.IContact
	ActivityModel models.ContactClientActivity
}

func NewContactClient() IContactClient {
	return &ContactClient{
		Model:         models.ClientContact{},
		ClientModel:   models.Client{},
		ContactModel:  models.Contact{},
		ActivityModel: models.ContactClientActivity{},
	}
}

func (p *ContactClient) Create(ctx context.Context, contactID int64, data *entities.Client) error {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	clientID, err := p.ClientModel.Create(ctx, data)
	if err != nil {
		return err
	}

	err = p.Model.Create(ctx, &entities.ClientContact{
		ClientID:  clientID,
		ContactID: contactID,
		CreatedBy: data.CreatedBy,
	})
	if err != nil {
		return err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	p.ActivityModel.Create(ctx, &entities.ContactClientActivity{
		CreatedBy: data.CreatedBy,
		ClientID:  clientID,
		ContactID: contactID,
		Type:      entities.ActivityCreated,
	})

	return nil
}

func (p *ContactClient) Add(ctx context.Context, data *entities.ClientContact) error {
	err := p.Model.Create(ctx, data)
	if err != nil {
		return err
	}

	p.ActivityModel.Create(ctx, &entities.ContactClientActivity{
		CreatedBy: data.CreatedBy,
		ClientID:  data.ClientID,
		ContactID: data.ContactID,
		Type:      entities.ActivityCreated,
	})
	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)
	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)

	return nil
}

func (p *ContactClient) AddBatch(ctx context.Context, data []entities.ClientContact) error {
	err := p.Model.CreateBatch(ctx, data)
	if err != nil {
		return err
	}

	for _, val := range data {
		p.ActivityModel.Create(ctx, &entities.ContactClientActivity{
			CreatedBy: val.CreatedBy,
			ClientID:  val.ClientID,
			ContactID: val.ContactID,
			Type:      entities.ActivityCreated,
		})

		p.ContactModel.UpdateLastActiveTime(ctx, val.ContactID)
		p.ClientModel.UpdateLastActiveTime(ctx, val.ClientID)
	}

	return nil
}

func (p *ContactClient) Delete(ctx context.Context, userID,contactID, clientID int64) error {

	err := p.Model.Delete(ctx, clientID, contactID)
	if err != nil {
		return err
	}
	p.ActivityModel.Create(ctx, &entities.ContactClientActivity{
		CreatedBy: userID,
		ClientID:  clientID,
		ContactID: contactID,
		Type:      entities.ActivityDeleted,
	})
	p.ClientModel.UpdateLastActiveTime(ctx, clientID)
	p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	return nil
}

func (p *ContactClient) List(ctx context.Context, contactID int64) ([]*entities.ClientContact, error) {
	filters := map[string]any{"contact_id": contactID}
	return p.Model.List(ctx, filters)
}
