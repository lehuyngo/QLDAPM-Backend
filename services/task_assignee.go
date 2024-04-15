package services

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type ITaskAssignee interface {
	Create(ctx context.Context, data *entities.TaskAssignee) error
	Delete(ctx context.Context, taskID, assigneeID int64) error
}

type TaskAssignee struct {
	Model models.ITaskAssignee
}

func NewTaskAssignee() ITaskAssignee {
	return &TaskAssignee{
		Model: models.TaskAssignee{},
	}
}

func (p *TaskAssignee) Create(ctx context.Context, data *entities.TaskAssignee) error {
	return p.Model.Create(ctx, data)
}

func (p *TaskAssignee) Delete(ctx context.Context, taskID, assigneeID int64) error {
	return p.Model.Delete(ctx, taskID, assigneeID)
}
