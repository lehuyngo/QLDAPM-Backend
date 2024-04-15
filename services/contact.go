package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContact interface {
	Create(ctx context.Context, data *entities.Contact) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Contact, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Contact, error)
	ReadByEmail(ctx context.Context, orgID int64, email string) (*entities.Contact, error)
	Update(ctx context.Context, data *entities.Contact) error
	Delete(ctx context.Context, userID, id int64) error
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Contact, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contact, error)
}

type Contact struct {
	Model         models.IContact
	ActivityModel models.ContactActivity
}

func NewContact() IContact {
	return &Contact{
		Model:         models.Contact{},
		ActivityModel: models.ContactActivity{},
	}
}

func (p *Contact) Create(ctx context.Context, data *entities.Contact) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	data.LastActiveTime = time.Now().UnixMilli()

	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}
	
	p.ActivityModel.Create(ctx, &entities.ContactActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		ContactID: data.ID,
		Type:      entities.ActivityCreated,
	})
	return result, err
}

func (p *Contact) Read(ctx context.Context, id int64) (*entities.Contact, error) {
	return p.Model.First(ctx, map[string]any{"id": id})
}

func (p *Contact) ReadByUUID(ctx context.Context, uuid string) (*entities.Contact, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Contact) ReadByEmail(ctx context.Context, orgID int64, email string) (*entities.Contact, error) {
	filters := map[string]any{"email": email,"organization_id": orgID}
	return p.Model.First(ctx, filters)
}

func (p *Contact) Update(ctx context.Context, data *entities.Contact) error {
	data.LastActiveTime = time.Now().UnixMilli()

	err := p.Model.Update(ctx, data)
	if err != nil {
		return err
	}

	typeActivity := entities.ActivityUpdated
	if data.Status == entities.InTrash {
		typeActivity = entities.ActivityDeleted
	}
	p.ActivityModel.Create(ctx, &entities.ContactActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.UpdatedBy,
		ContactID: data.ID,
		Type:      typeActivity,
	})

	return nil
}

func (p *Contact) Delete(ctx context.Context, userID,id int64) error {
	error := p.Model.Delete(ctx, id)
	if error != nil {
		return error
	}
	p.ActivityModel.Create(ctx, &entities.ContactActivity{
		UUID:      uuid.NewString(),
		CreatedBy: userID,
		ContactID: id,
		Type:      entities.ActivityDeleted,
	})
	return nil
}

func (p *Contact) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Contact, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *Contact) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contact, error) {
	return p.Model.ListByUUIDs(ctx, uuids)
}
