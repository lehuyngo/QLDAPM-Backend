package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContactProject interface {
	Create(ctx context.Context, contactID int64, data *entities.Project) error
	Add(ctx context.Context, data *entities.ContactProject) error
	AddBatch(ctx context.Context, data []entities.ContactProject) error
	Delete(ctx context.Context, userID, contactID, projectID int64) error
	List(ctx context.Context, contactID int64) ([]*entities.ContactProject, error)
}

type ContactProject struct {
	Model           models.IContactProject
	ContactModel    models.IContact
	ProjectModel    models.IProject
	ActivitiesModel models.IContactProjectActivity
}

func NewContactProject() IContactProject {
	return &ContactProject{
		Model:           models.ContactProject{},
		ContactModel:    models.Contact{},
		ProjectModel:    models.Project{},
		ActivitiesModel: models.ContactProjectActivity{},
	}
}

func (p *ContactProject) Create(ctx context.Context, contactID int64, data *entities.Project) error {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	data.LastActiveTime = time.Now().UnixMilli()
	projectID, err := p.ProjectModel.Create(ctx, data)
	if err != nil {
		return err
	}

	err = p.Model.Create(ctx, &entities.ContactProject{
		ContactID: contactID,
		ProjectID: projectID,
		CreatedBy: data.CreatedBy,
	})
	if err != nil {
		return err
	}
	p.ActivitiesModel.Create(ctx, &entities.ContactProjectActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		ProjectID: projectID,
		ContactID: contactID,
		Type:      entities.ActivityCreated,
	})

	p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	return nil
}

func (p *ContactProject) Add(ctx context.Context, data *entities.ContactProject) error {
	err := p.Model.Create(ctx, data)
	if err != nil {
		return err
	}
	p.ActivitiesModel.Create(ctx, &entities.ContactProjectActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		ProjectID: data.ProjectID,
		ContactID: data.ContactID,
		Type:      entities.ActivityCreated,
	})
	p.ContactModel.UpdateLastActiveTime(ctx, data.ContactID)
	p.ProjectModel.UpdateLastActiveTime(ctx, data.ProjectID)
	return nil
}

func (p *ContactProject) AddBatch(ctx context.Context, data []entities.ContactProject) error {
	err := p.Model.CreateBatch(ctx, data)
	if err != nil {
		return err
	}

	for _, val := range data {
		p.ActivitiesModel.Create(ctx, &entities.ContactProjectActivity{
			UUID:      uuid.NewString(),
			CreatedBy: val.CreatedBy,
			ProjectID: val.ProjectID,
			ContactID: val.ContactID,
			Type:      entities.ActivityCreated,
		})
		p.ContactModel.UpdateLastActiveTime(ctx, val.ContactID)
		p.ProjectModel.UpdateLastActiveTime(ctx, val.ProjectID)
	}
	return nil
}

func (p *ContactProject) Delete(ctx context.Context, userID, contactID, projectID int64) error {
	err := p.Model.Delete(ctx, contactID, projectID)
	if err != nil {
		return err
	}
	p.ActivitiesModel.Create(ctx, &entities.ContactProjectActivity{
		UUID:      uuid.NewString(),
		CreatedBy: userID,
		ProjectID: projectID,
		ContactID: contactID,
		Type:      entities.ActivityDeleted,
	})
	p.ContactModel.UpdateLastActiveTime(ctx, contactID)
	p.ProjectModel.UpdateLastActiveTime(ctx, projectID)
	return nil
}

func (p *ContactProject) List(ctx context.Context, contactID int64) ([]*entities.ContactProject, error) {
	filters := map[string]any{"contact_id": contactID}
	return p.Model.List(ctx, filters)
}
