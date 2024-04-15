package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IMeetingNoteEditor interface {
	Create(ctx context.Context, data *entities.MeetingNoteEditor) (int64, error)
}

type MeetingNoteEditor struct {
	Model models.IMeetingNoteEditor
}

func NewMeetingNoteEditor() IMeetingNoteEditor {
	return &MeetingNoteEditor{
		Model: models.MeetingNoteEditor{},
	}
}
func (p *MeetingNoteEditor) Create(ctx context.Context, data *entities.MeetingNoteEditor) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}
