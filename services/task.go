package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type ITask interface {
	Create(ctx context.Context, data *entities.Task) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Task, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Task, error)
	Update(ctx context.Context, data *entities.Task) error
	Delete(ctx context.Context, id int64) error
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Task, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*entities.Task, error)
}

type Task struct {
	Model models.ITask
}

func NewTask() ITask {
	return &Task{
		Model: models.Task{},
	}
}

func (p *Task) Create(ctx context.Context, data *entities.Task) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *Task) Read(ctx context.Context, id int64) (*entities.Task, error) {
	return p.Model.Read(ctx, id)
}

func (p *Task) ReadByUUID(ctx context.Context, uuid string) (*entities.Task, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Task) Update(ctx context.Context, data *entities.Task) error {
	return p.Model.Update(ctx, data)
}

func (p *Task) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *Task) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Task, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *Task) ListByProjectID(ctx context.Context, projectID int64) ([]*entities.Task, error) {
	filters := map[string]any{"project_id": projectID}
	return p.Model.List(ctx, filters)
}
