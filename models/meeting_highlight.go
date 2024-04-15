package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/gorm"
)

type IMeetingHighlight interface {
	Create(ctx context.Context, data *entities.MeetingHighlight) (int64, error)
	CreateBatch(ctx context.Context, data []*entities.MeetingHighlight) error
	List(ctx context.Context, filters map[string]any) ([]*entities.MeetingHighlight, error)
	ListByMeetingNotes(ctx context.Context, meetingNoteIDs []int64) ([]*entities.MeetingHighlight, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingHighlight, error)
	Delete(ctx context.Context, id int64) error
	DeleteBatch(ctx context.Context, uuids []string, meetingNoteID int64) error
	DeleteWithCondition(ctx context.Context, filters map[string]any) error
	First(ctx context.Context, filters map[string]any) (*entities.MeetingHighlight, error)
}

type MeetingHighlight struct {
}

func (MeetingHighlight) Create(ctx context.Context, data *entities.MeetingHighlight) (int64, error) {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return data.ID, err
}

func (MeetingHighlight) CreateBatch(ctx context.Context, data []*entities.MeetingHighlight) error {
	err := clients.MySQLClient.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(data).Error
	})

	return err
}

func (MeetingHighlight) List(ctx context.Context, filters map[string]any) ([]*entities.MeetingHighlight, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := []*entities.MeetingHighlight{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("MeetingNote").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingHighlight) ListByMeetingNotes(ctx context.Context, meetingNoteIDs []int64) ([]*entities.MeetingHighlight, error) {

	result := []*entities.MeetingHighlight{}

	db := clients.MySQLClient.WithContext(ctx)

	err := db.Preload("MeetingNote").Preload("Creator").Where("meeting_note_id IN ?", meetingNoteIDs).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingHighlight) ReadByUUID(ctx context.Context, uuid string) (*entities.MeetingHighlight, error) {
	result := &entities.MeetingHighlight{}
	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.MeetingHighlight{})
	err := db.Preload("MeetingNote.Project").Where("uuid = ?", uuid).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingHighlight) Delete(ctx context.Context, id int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Delete(&entities.MeetingHighlight{}, id).Error
}

func (MeetingHighlight) DeleteWithCondition(ctx context.Context, filters map[string]any) error {
	db := clients.MySQLClient

	query := db.WithContext(ctx)
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Delete(&entities.MeetingHighlight{}).Error; err != nil {
		return err
	}

	return nil
}

func (MeetingHighlight) First(ctx context.Context, filters map[string]any) (*entities.MeetingHighlight, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}

	result := &entities.MeetingHighlight{}

	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName())
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.
		WithContext(ctx).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (MeetingHighlight) DeleteBatch(ctx context.Context, uuids []string, meetingNoteID int64) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Where("uuid IN ? AND meeting_note_id = ?", uuids, meetingNoteID).Delete(&entities.MeetingHighlight{}).Error
}
