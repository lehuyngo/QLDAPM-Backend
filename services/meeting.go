package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IMeeting interface {
	Create(ctx context.Context, data *entities.Meeting) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Meeting, error)
	List(ctx context.Context, projectId int64) ([]*entities.Meeting, error)
}

type Meeting struct {
	Model        models.IMeeting
	ProjectModel models.IProject
}

func NewMeeting() IMeeting {
	return &Meeting{
		Model:        models.Meeting{},
		ProjectModel: models.Project{},
	}
}

func (p *Meeting) Create(ctx context.Context, data *entities.Meeting) (int64, error) {
	data.UUID = uuid.NewString()
	data.LastActiveTime = time.Now().UnixMilli()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	// update last active time of project
	p.ProjectModel.UpdateLastActiveTime(ctx, data.ProjectID)

	return result, nil
}

func (p *Meeting) ReadByUUID(ctx context.Context, uuid string) (*entities.Meeting, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Meeting) List(ctx context.Context, projectID int64) ([]*entities.Meeting, error) {
	return p.Model.List(ctx, map[string]any{"project_id": projectID})
}
