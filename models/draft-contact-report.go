package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IDraftContactReport interface {
	ListAll(ctx context.Context, filters map[string]any) ([]*entities.DraftContact, error)
}

type DraftContactReport struct {
}

func (DraftContactReport) ListAll(ctx context.Context, filters map[string]any) ([]*entities.DraftContact, error) {
	result := []*entities.DraftContact{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Unscoped().Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
