package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IStaticFile interface {
	Create(ctx context.Context, data *entities.StaticFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.StaticFile, error)
}

type StaticFile struct {
	Model models.IStaticFile
}

func NewStaticFile() IStaticFile {
	return &StaticFile{
		Model: models.StaticFile{},
	}
}

func (p *StaticFile) Create(ctx context.Context, data *entities.StaticFile) (string, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *StaticFile) Read(ctx context.Context, uuid string) (*entities.StaticFile, error) {
	return p.Model.Read(ctx, uuid)
}
