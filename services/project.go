package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IProject interface {
	Create(ctx context.Context, data *entities.Project) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Project, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Project, error)
	Update(ctx context.Context, data *entities.Project) error
	Delete(ctx context.Context, id int64) error
	ListByClientID(ctx context.Context, clientID int64) ([]*entities.Project, error)
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Project, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Project, error)
}

type Project struct {
	Model                     models.IProject
	ClientProjectAcivityModel models.IClientProjectActivity
}

func NewProject() IProject {
	return &Project{
		Model:                     models.Project{},
		ClientProjectAcivityModel: models.ClientProjectActivity{},
	}
}

func (p *Project) Create(ctx context.Context, data *entities.Project) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	data.LastActiveTime = time.Now().UnixMilli()

	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	if data.ClientID != 0 {
		p.ClientProjectAcivityModel.Create(ctx, &entities.ClientProjectActivity{
			ClientID:  data.ClientID,
			ProjectID: result,
			Type:      entities.ActivityCreated,
			CreatedBy: data.CreatedBy,
		})
	}
	
	return result, nil
}

func (p *Project) Read(ctx context.Context, id int64) (*entities.Project, error) {
	return p.Model.First(ctx, map[string]any{"id": id})
}

func (p *Project) ReadByUUID(ctx context.Context, uuid string) (*entities.Project, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Project) Update(ctx context.Context, data *entities.Project) error {
	data.LastActiveTime = time.Now().UnixMilli()
	return p.Model.Update(ctx, data)
}

func (p *Project) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *Project) ListByClientID(ctx context.Context, clientID int64) ([]*entities.Project, error) {
	filters := map[string]any{"client_id": clientID}
	return p.Model.List(ctx, filters)
}

func (p *Project) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Project, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *Project) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Project, error) {
	return p.Model.ListByUUIDs(ctx, uuids)
}
