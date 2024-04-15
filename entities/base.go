package entities

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType int

const (
	ActivityCreated ActivityType = 0
	ActivityDeleted ActivityType = 1
	ActivityUpdated ActivityType = 2
)

type TaskStatus int

const (
	TaskStatusToDo    TaskStatus = 1
	TaskStatusDoing   TaskStatus = 2
	TaskStatusTesting TaskStatus = 3
	TaskStatusDone    TaskStatus = 4
)

type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 1
	TaskPriorityMedium TaskPriority = 2
	TaskPriorityHigh   TaskPriority = 3
)

type TaskLabel int

const (
	TaskLabelTask        TaskLabel = 1
	TaskLabelFeedback    TaskLabel = 2
	TaskLabelImprovement TaskLabel = 3
	TaskLabelBug         TaskLabel = 4
)

type Base struct {
	UUID      string         `gorm:"column:uuid;omitempty"`
	CreatedBy int64          `gorm:"column:created_by;omitempty"`
	Creator   *User          `gorm:"references:created_by;foreignKey:id;omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;omitempty"`
	UpdatedBy int64          `gorm:"column:updated_by;omitempty"`
	Updater   *User          `gorm:"references:updated_by;foreignKey:id;omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `swaggertype:"string" gorm:"index;column:deleted_at"`
}

func (b *Base) GetCreator() *User {
	if b == nil {
		return nil
	}

	return b.Creator
}

func (b *Base) GetCreatedAt() time.Time {
	if b == nil {
		return time.Now()
	}

	return b.CreatedAt
}

func (b *Base) GetUpdater() *User {
	if b == nil {
		return nil
	}

	return b.Updater
}

func (b *Base) GetUpdatedAt() time.Time {
	if b == nil {
		return time.Now()
	}

	return b.UpdatedAt
}

func (b *Base) GetUUID() string {
	if b == nil {
		return ""
	}

	return b.UUID
}

type MailStatus int

const (
	MailFailed 		MailStatus = -1
	MailWaiting 	MailStatus = 1
	MailProcessing 	MailStatus = 2
	MailSuccess		MailStatus = 3
)
