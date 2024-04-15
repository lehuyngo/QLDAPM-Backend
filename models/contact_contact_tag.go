package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IContactContactTag interface {
	First(ctx context.Context, tagID int64) (*entities.ContactContactTag, error)
	Create(ctx context.Context, data *entities.ContactContactTag) error
	Delete(ctx context.Context, contactID, tagID int64) error
}

type ContactContactTag struct {
}

func (ContactContactTag) First(ctx context.Context, tagID int64) (*entities.ContactContactTag, error) {
	result := &entities.ContactContactTag{}
	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.ContactContactTag{})
	err := db.Where("tag_id = ?", tagID).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ContactContactTag) Create(ctx context.Context, data *entities.ContactContactTag) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(data).Error
}

func (ContactContactTag) Delete(ctx context.Context, contactID, tagID int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("contact_id = ? AND tag_id = ?", contactID, tagID).Delete(&entities.ContactContactTag{}).Error
}
