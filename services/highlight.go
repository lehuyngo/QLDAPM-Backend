package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IMeetingHighlight interface {
	Create(ctx context.Context, data *entities.MeetingHighlight) (int64, error)
	CreateBatch(ctx context.Context, data []*entities.MeetingHighlight) error
	ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingHighlight, error)
	Delete(ctx context.Context, id int64) error
	DeleteBatch(ctx context.Context, uuids []string, meetingNoteID int64) error
	ListByMeetingNotes(ctx context.Context, meetingNoteIDs []int64) ([]*entities.MeetingHighlight, error)
	First(ctx context.Context, title string, meetingNoteID int64) (*entities.MeetingHighlight, error)
}

type MeetingHighlight struct {
	Model models.IMeetingHighlight
}

func NewMeetingHighlight() IMeetingHighlight {
	return &MeetingHighlight{
		Model: models.MeetingHighlight{},
	}
}

func (p *MeetingHighlight) Create(ctx context.Context, data *entities.MeetingHighlight) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *MeetingHighlight) CreateBatch(ctx context.Context, data []*entities.MeetingHighlight) error {
	return p.Model.CreateBatch(ctx, data)
}

func (p *MeetingHighlight) ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingHighlight, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *MeetingHighlight) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *MeetingHighlight) ListByMeetingNotes(ctx context.Context, meetingNoteIDs []int64) ([]*entities.MeetingHighlight, error) {
	return p.Model.ListByMeetingNotes(ctx, meetingNoteIDs)
}
func (p *MeetingHighlight) First(ctx context.Context, title string, meetingNoteID int64) (*entities.MeetingHighlight, error) {
	filters := map[string]any{"title": title, "meeting_note_id": meetingNoteID}
	return p.Model.First(ctx, filters)
}
func (p *MeetingHighlight) DeleteBatch(ctx context.Context, uuids []string, meetingNoteID int64) error {
	return p.Model.DeleteBatch(ctx, uuids, meetingNoteID)
}
