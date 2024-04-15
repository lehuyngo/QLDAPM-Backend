package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type ITaskAssignee interface {
	Create(ctx context.Context, data *entities.TaskAssignee) error
	Delete(ctx context.Context, taskID, assigneeID int64) error
}

type TaskAssignee struct {
}

func (TaskAssignee) Create(ctx context.Context, data *entities.TaskAssignee) error {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})
	return err
}

func (TaskAssignee) Delete(ctx context.Context, taskID, assigneeID int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("task_id = ? AND assignee_id = ?", taskID, assigneeID).Delete(&entities.TaskAssignee{}).Error
}
