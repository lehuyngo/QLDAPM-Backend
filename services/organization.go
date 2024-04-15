package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IOrganization interface {
	Create(ctx context.Context, data *entities.Organization) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Organization, error)
	First(ctx context.Context) (*entities.Organization, error)
	Update(ctx context.Context, data *entities.Organization) error
}

type Organization struct {
	Model models.IOrganization
}

func NewOrganization() IOrganization {
	return &Organization{
		Model: models.Organization{},
	}
}

func (p *Organization) Create(ctx context.Context, data *entities.Organization) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *Organization) Read(ctx context.Context, id int64) (*entities.Organization, error) {
	return p.Model.Read(ctx, id)
}

func (p *Organization) First(ctx context.Context) (*entities.Organization, error) {
	return p.Model.First(ctx)
}

func (p *Organization) Update(ctx context.Context, data *entities.Organization) error {
	return p.Model.Update(ctx, data)
}
