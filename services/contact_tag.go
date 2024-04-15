package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gorm.io/gorm"
)

type IContactTag interface {
	Create(ctx context.Context, data *entities.ContactTag) (int64, error)
	Add(ctx context.Context, data *entities.ContactContactTag) error
	Read(ctx context.Context, id int64) (*entities.ContactTag, error)
	ReadByUUID(ctx context.Context, orgID int64, uuid string) (*entities.ContactTag, error)
	ReadByName(ctx context.Context, orgID int64, name string) (*entities.ContactTag, error)
	Update(ctx context.Context, data *entities.ContactTag) (error)
	List(ctx context.Context, orgID int64) ([]*entities.ContactTag, error)
	Delete(ctx context.Context, user *entities.User, contactID, tagID int64) error
}

type ContactTag struct {
	Model         models.IContactTag
	JoinModel     models.IContactContactTag
	ContactModel  models.IContact
	ActivityModel models.IContactTagActivity
}

func NewContactTag() IContactTag {
	return &ContactTag{
		Model:         models.ContactTag{},
		JoinModel:     models.ContactContactTag{},
		ContactModel:  models.Contact{},
		ActivityModel: models.ContactTagActivity{},
	}
}

func (p *ContactTag) Create(ctx context.Context, data *entities.ContactTag) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	var contactID int64
	if len(data.Contacts) > 0 {
		contactID = data.Contacts[0].ContactID
		p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	}

	// add contact tags activity
	p.ActivityModel.Create(ctx, &entities.ContactTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		TagID:     data.ID,
		Type:      entities.ActivityCreated,
		ContactID: contactID,
	})

	return result, nil
}

func (p *ContactTag) Add(ctx context.Context, data *entities.ContactContactTag) error {
	err := p.JoinModel.Create(ctx, data)
	if err != nil {
		return err
	}

	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)

	// add contact tags activity
	p.ActivityModel.Create(ctx, &entities.ContactTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		TagID:     data.TagID,
		Type:      entities.ActivityCreated,
		ContactID: data.ContactID,
	})

	return nil
}

func (p *ContactTag) Read(ctx context.Context, id int64) (*entities.ContactTag, error) {
	return p.Model.First(ctx, map[string]any{"id": id})
}

func (p *ContactTag) ReadByUUID(ctx context.Context, orgID int64, uuid string) (*entities.ContactTag, error) {
	return p.Model.First(ctx, map[string]any{"organization_id": orgID, "uuid": uuid})
}

func (p *ContactTag) ReadByName(ctx context.Context, orgID int64, name string) (*entities.ContactTag, error) {
	return p.Model.First(ctx, map[string]any{"organization_id": orgID, "name": name})
}

func (p *ContactTag) Update(ctx context.Context, data *entities.ContactTag) (error) {
	return p.Model.Update(ctx, data)
}

func (p *ContactTag) List(ctx context.Context, orgID int64) ([]*entities.ContactTag, error) {
	return p.Model.List(ctx, map[string]any{"organization_id": orgID})
}

func (p *ContactTag) Delete(ctx context.Context, user *entities.User, contactID, tagID int64) error {
	err := p.JoinModel.Delete(ctx, contactID, tagID)
	if err != nil {
		return err
	}

	_, err = p.JoinModel.First(ctx, tagID)
	if err == nil {
		return nil
	}

	// remove untagged tag
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.Model.Delete(ctx, tagID)
	}

	p.ContactModel.UpdateLastActiveTime(ctx, contactID)

	// add contact tags activity
	p.ActivityModel.Create(ctx, &entities.ContactTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: user.ID,
		TagID:     tagID,
		Type:      entities.ActivityDeleted,
		ContactID: contactID,
	})

	return nil
}
