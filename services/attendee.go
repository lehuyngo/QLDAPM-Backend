package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IAttendee interface {
	Create(ctx context.Context, data *entities.Attendee) (int64, error)
}

type Attendee struct {
	Model models.IAttendee
}

func NewAttendee() IAttendee {
	return &Attendee{
		Model: models.Attendee{},
	}
}
func (p *Attendee) Create(ctx context.Context, data *entities.Attendee) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return result, nil

}
