package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IContributor interface {
	Create(ctx context.Context, data *entities.Contributor) (int64, error)
	CreateBatch(ctx context.Context, data []*entities.Contributor) error
	ReadByUUID(ctx context.Context, uuid string) (*entities.Contributor, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contributor, error)
	Delete(ctx context.Context, id int64) error
	DeleteBatch(ctx context.Context, uuids []string) error
}

type Contributor struct {
	Model models.IContributor
}

func NewContributor() IContributor {
	return &Contributor{
		Model: models.Contributor{},
	}
}
func (p *Contributor) Create(ctx context.Context, data *entities.Contributor) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *Contributor) CreateBatch(ctx context.Context, data []*entities.Contributor) error {
	return p.Model.CreateBatch(ctx, data)
}

func (p *Contributor) ReadByUUID(ctx context.Context, uuid string) (*entities.Contributor, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Contributor) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Contributor, error) {
	return p.Model.ListByUUIDs(ctx, uuids)
}

func (p *Contributor) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *Contributor) DeleteBatch(ctx context.Context, uuids []string) error {
	return p.Model.DeleteBatch(ctx, uuids)
}
