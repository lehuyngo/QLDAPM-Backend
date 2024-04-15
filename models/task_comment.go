package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type ITaskComment interface {
	Create(ctx context.Context, data *entities.TaskComment) (int64, error)
	List(ctx context.Context, filters map[string]any) ([]*entities.TaskComment, error)
}
type TaskComment struct {
}

func (p TaskComment) Create(ctx context.Context, data *entities.TaskComment) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(data).Error
	})
	return data.ID, err
}
func (p TaskComment) List(ctx context.Context, filters map[string]any) ([]*entities.TaskComment, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	result := []*entities.TaskComment{}
	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Creator").Find(&result).Error
	return result, err
}
