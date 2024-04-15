package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type ITaskAttachFile interface {
	Create(ctx context.Context, data *entities.TaskAttachFile) (string, error)
	Read(ctx context.Context, uuid string) (*entities.TaskAttachFile, error)
	Delete(ctx context.Context, taskID, id int64) error
}

type TaskAttachFile struct {
	Model models.ITaskAttachFile
}

func NewTaskAttachFile() ITaskAttachFile {
	return &TaskAttachFile{
		Model: models.TaskAttachFile{},
	}
}

func (p *TaskAttachFile) Create(ctx context.Context, data *entities.TaskAttachFile) (string, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *TaskAttachFile) Read(ctx context.Context, uuid string) (*entities.TaskAttachFile, error) {
	return p.Model.Read(ctx, uuid)
}

func (p *TaskAttachFile) Delete(ctx context.Context, taskID, id int64) error {
	return p.Model.Delete(ctx, id)
}
