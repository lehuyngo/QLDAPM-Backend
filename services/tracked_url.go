package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type ITrackedURL interface {
	Create(ctx context.Context, data *entities.TrackedURL) (int64, error)
	Read(ctx context.Context, uuid string) (*entities.TrackedURL, error)
	ReadByCode(ctx context.Context, code string) (*entities.TrackedURL, error)
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.TrackedURL, error)
	SetRead(ctx context.Context, id int64) error
	Update(ctx context.Context, data *entities.TrackedURL) error
	Delete(ctx context.Context, id int64) error
}

type TrackedURL struct {
	Model models.ITrackedURL
}

func NewTrackedURL() ITrackedURL {
	return &TrackedURL{
		Model: models.TrackedURL{},
	}
}

func (p *TrackedURL) Create(ctx context.Context, data *entities.TrackedURL) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.New
	return p.Model.Create(ctx, data)
}

func (p *TrackedURL) Read(ctx context.Context, uuid string) (*entities.TrackedURL, error) {
	result, err := p.Model.ReadByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	result.Status =  entities.Read
	if result.Status != entities.Read {
		p.Model.UpdateStatus(ctx, result.ID, entities.Read)
	}
	
	return result, nil
}

func (p *TrackedURL) ReadByCode(ctx context.Context, code string) (*entities.TrackedURL, error) {
	filters := map[string]any{"code": code}
	return p.Model.First(ctx, filters)
}

func (p *TrackedURL) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.TrackedURL, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *TrackedURL) SetRead(ctx context.Context, id int64) error {
	return p.Model.UpdateStatus(ctx, id, entities.Read)
}

func (p *TrackedURL) Update(ctx context.Context, data *entities.TrackedURL) error {
	return p.Model.Update(ctx, data)
}

func (p *TrackedURL) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}
