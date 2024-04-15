package services

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IDraftContactReport interface {
	GetReport(ctx context.Context, orgID int64) (*entities.DraftContactReport, error)
}

type DraftContactReport struct {
	Model models.IDraftContactReport
}

func NewDraftContactReport() IDraftContactReport {
	return &DraftContactReport{
		Model: models.DraftContactReport{},
	}
}
func (p *DraftContactReport) GetReport(ctx context.Context, orgID int64) (*entities.DraftContactReport, error) {
	filters := map[string]interface{}{"organization_id": orgID}
	data, err := p.Model.ListAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	total := len(data)
	resolved := 0
	processing := 0

	for _, draftContact := range data {
		if draftContact.ContactID != 0 {
			resolved++
		} else if !draftContact.DeletedAt.Valid {
			processing++
		}
	}

	deleted := total - processing - resolved
	result := &entities.DraftContactReport{
		Total:      total,
		Resolved:   resolved,
		Processing: processing,
		Deleted:    deleted,
	}
	return result, nil
}
