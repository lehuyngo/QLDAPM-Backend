package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IMeetingNote interface {
	Create(ctx context.Context, data *entities.MeetingNote) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingNote, error)
	List(ctx context.Context, projectID int64) ([]*entities.MeetingNote, error)
	ListByMeetingID(ctx context.Context, meetingID int64) ([]*entities.MeetingNote, error)
	Update(ctx context.Context, data *entities.MeetingNote) error
	Delete(ctx context.Context, id int64, projectID int64, meetingID int64) error
}

type MeetingNote struct {
	Model          models.IMeetingNote
	HighlightModel models.IMeetingHighlight
	ProjectModel   models.IProject
	MeetingModel   models.IMeeting
}

func NewMeetingNote() IMeetingNote {
	return &MeetingNote{
		Model:          models.MeetingNote{},
		HighlightModel: models.MeetingHighlight{},
		ProjectModel:   models.Project{},
		MeetingModel:   models.Meeting{},
	}
}

func (p *MeetingNote) Create(ctx context.Context, data *entities.MeetingNote) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	// update last active time of project
	p.ProjectModel.UpdateLastActiveTime(ctx, data.ProjectID)

	// update last active time of meeting
	p.MeetingModel.UpdateLastActiveTime(ctx, data.MeetingID)

	return result, nil
}

func (p *MeetingNote) ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingNote, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *MeetingNote) List(ctx context.Context, projectID int64) ([]*entities.MeetingNote, error) {
	return p.Model.List(ctx, map[string]any{"project_id": projectID})
}

func (p *MeetingNote) ListByMeetingID(ctx context.Context, meetingID int64) ([]*entities.MeetingNote, error) {
	return p.Model.List(ctx, map[string]any{"meeting_id": meetingID})
}

func (p *MeetingNote) Update(ctx context.Context, data *entities.MeetingNote) error {
	err := p.Model.Update(ctx, data)
	if err != nil {
		return err
	}

	// update last active time of project
	p.ProjectModel.UpdateLastActiveTime(ctx, data.ProjectID)

	// update last active time of meeting
	p.MeetingModel.UpdateLastActiveTime(ctx, data.MeetingID)

	return nil
}

func (p *MeetingNote) Delete(ctx context.Context, id int64, projectID int64, meetingID int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = p.HighlightModel.DeleteWithCondition(ctx, map[string]any{"meeting_note_id": id})

	// update last active time of project
	p.ProjectModel.UpdateLastActiveTime(ctx, projectID)

	// update last active time of meeting
	p.MeetingModel.UpdateLastActiveTime(ctx, meetingID)

	return err
}
