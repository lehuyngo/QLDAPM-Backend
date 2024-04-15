package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IBatchMail interface {
	Create(ctx context.Context, data *entities.BatchMail) (int64, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.BatchMail, error)
	UpdateStatus(ctx context.Context, id int64, status entities.MailStatus) error
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.BatchMail, error)
}

type BatchMail struct {
	Model models.IBatchMail
}

func NewBatchMail() IBatchMail {
	return &BatchMail{
		Model: models.BatchMail{},
	}
}

func (p *BatchMail) Create(ctx context.Context, data *entities.BatchMail) (int64, error) {
	data.UUID = uuid.NewString()
	return p.Model.Create(ctx, data)
}

func (p *BatchMail) ReadByUUID(ctx context.Context, uuid string) (*entities.BatchMail, error) {
	result, err := p.Model.ReadByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	p.updateLastStatus(ctx, result)
	return result, nil
}

func (p *BatchMail) UpdateStatus(ctx context.Context, id int64, status entities.MailStatus) error {
	return p.Model.UpdateStatus(ctx, id, status)
}

func (p *BatchMail) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.BatchMail, error) {
	filters := map[string]any{"organization_id": orgID}
	result, err := p.Model.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	for _, val := range result {
		p.updateLastStatus(ctx, val)
	}

	return result, nil
}

func (p *BatchMail) updateLastStatus(ctx context.Context, val *entities.BatchMail) (error) {
	if (val.Status == entities.MailSuccess) || (val.Status == entities.MailFailed) {
		return nil
	}

	successCounter := 0
	for _, receiver := range val.Receivers {
		if receiver.Status == entities.MailFailed {
			val.Status = entities.MailFailed
			return p.UpdateStatus(ctx, val.ID, entities.MailFailed)
		}

		if receiver.Status == entities.MailSuccess {
			successCounter = successCounter + 1
			continue
		}
	}

	if (successCounter > 0) && (successCounter == len(val.Receivers)) {
		val.Status = entities.MailSuccess
		return p.UpdateStatus(ctx, val.ID, entities.MailSuccess)
	}

	return nil
}
