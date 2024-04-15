package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactNote interface {
	Create(ctx context.Context, data *entities.ContactNote) (int64, error)
	Read(ctx context.Context, id int64) (*entities.ContactNote, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.ContactNote, error)
	Update(ctx context.Context, data *entities.ContactNote) error
	MoveToTrash(ctx context.Context, id int64) error
	Delete(ctx context.Context, user *entities.User, contactID, id int64) error
	ListByContactID(ctx context.Context, contactID int64) ([]*entities.ContactNote, error)
}

type ContactNote struct {
	Model         models.IContactNote
	ContactModel  models.Contact
	ActivityModel models.ContactNoteActivity
}

func NewContactNote() IContactNote {
	return &ContactNote{
		Model:         models.ContactNote{},
		ContactModel:  models.Contact{},
		ActivityModel: models.ContactNoteActivity{},
	}
}

func (p *ContactNote) Create(ctx context.Context, data *entities.ContactNote) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)

	// add contact notes activity
	p.ActivityModel.Create(ctx, &entities.ContactNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		NoteID:    data.ID,
		Type:      entities.ActivityCreated,
		ContactID: data.ContactID,
	})

	return result, nil
}

func (p *ContactNote) Read(ctx context.Context, id int64) (*entities.ContactNote, error) {
	return p.Model.Read(ctx, id)
}

func (p *ContactNote) ReadByUUID(ctx context.Context, uuid string) (*entities.ContactNote, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *ContactNote) Update(ctx context.Context, data *entities.ContactNote) error {
	err := p.Model.Update(ctx, data)
	if err != nil {
		return err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)

	// add contact notes activity
	p.ActivityModel.Create(ctx, &entities.ContactNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.UpdatedBy,
		NoteID:    data.ID,
		Type:      entities.ActivityUpdated,
		ContactID: data.ContactID,
	})

	return nil
}

func (p *ContactNote) MoveToTrash(ctx context.Context, id int64) error {
	err := p.Model.UpdateField(ctx, id, "status", entities.InTrash)
	if err != nil {
		return err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, id)
	return nil
}

func (p *ContactNote) Delete(ctx context.Context, user *entities.User, contactID, id int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, contactID)

	// add contact notes activity
	p.ActivityModel.Create(ctx, &entities.ContactNoteActivity{
		UUID:      uuid.NewString(),
		CreatedBy: user.ID,
		NoteID:    id,
		Type:      entities.ActivityDeleted,
		ContactID: contactID,
	})

	return nil
}

func (p *ContactNote) ListByContactID(ctx context.Context, contactID int64) ([]*entities.ContactNote, error) {
	filters := map[string]any{"contact_id": contactID}
	return p.Model.List(ctx, filters)
}
