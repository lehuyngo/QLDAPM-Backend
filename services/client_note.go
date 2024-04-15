package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientNote interface {
	Create(ctx context.Context, data *entities.ClientNote) (int64, error)
	Read(ctx context.Context, id int64) (*entities.ClientNote, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.ClientNote, error)
	Update(ctx context.Context, data *entities.ClientNote) error
	MoveToTrash(ctx context.Context, clientID, id int64) error
	Delete(ctx context.Context, user *entities.User, clientID, id int64) error
	ListByClientID(ctx context.Context, clientID int64) ([]*entities.ClientNote, error)
}

type ClientNote struct {
	Model         models.IClientNote
	ClientModel   models.IClient
	ActivityModel models.IClientNoteActivity
}

func NewClientNote() IClientNote {
	return &ClientNote{
		Model:         models.ClientNote{},
		ClientModel:   models.Client{},
		ActivityModel: models.ClientNoteActivity{},
	}
}

func (p *ClientNote) Create(ctx context.Context, data *entities.ClientNote) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)

	// add client notes activity
	p.ActivityModel.Create(ctx, &entities.ClientNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		NoteID:    data.ID,
		Type:      entities.ActivityCreated,
		ClientID:  data.ClientID,
	})

	return result, nil
}

func (p *ClientNote) Read(ctx context.Context, id int64) (*entities.ClientNote, error) {
	return p.Model.Read(ctx, id)
}

func (p *ClientNote) ReadByUUID(ctx context.Context, uuid string) (*entities.ClientNote, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *ClientNote) Update(ctx context.Context, data *entities.ClientNote) error {
	err := p.Model.Update(ctx, data)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)

	// add client notes activity
	p.ActivityModel.Create(ctx, &entities.ClientNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.UpdatedBy,
		NoteID:    data.ID,
		Type:      entities.ActivityUpdated,
		ClientID:  data.ClientID,
	})

	return nil
}

func (p *ClientNote) MoveToTrash(ctx context.Context, clientID, id int64) error {
	err := p.Model.UpdateField(ctx, id, "status", entities.InTrash)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)
	return nil
}

func (p *ClientNote) Delete(ctx context.Context, user *entities.User, clientID, id int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)

	// add contact notes activity
	p.ActivityModel.Create(ctx, &entities.ClientNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: user.ID,
		NoteID:    id,
		Type:      entities.ActivityDeleted,
		ClientID:  clientID,
	})

	return nil
}

func (p *ClientNote) ListByClientID(ctx context.Context, clientID int64) ([]*entities.ClientNote, error) {
	filters := map[string]any{"client_id": clientID}
	return p.Model.List(ctx, filters)
}
