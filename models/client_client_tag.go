package models

import (
	"context"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IClientClientTag interface {
	First(ctx context.Context, tagID int64) (*entities.ClientClientTag, error)
	Create(ctx context.Context, data *entities.ClientClientTag) error
	Delete(ctx context.Context, clientID, tagID int64) error
}

type ClientClientTag struct {
}

func (ClientClientTag) First(ctx context.Context, tagID int64) (*entities.ClientClientTag, error) {
	result := &entities.ClientClientTag{}
	db := clients.MySQLClient.WithContext(ctx).Table(result.TableName()).Model(&entities.ClientClientTag{})
	err := db.Where("tag_id = ?", tagID).First(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ClientClientTag) Create(ctx context.Context, data *entities.ClientClientTag) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(data).Error
}

func (ClientClientTag) Delete(ctx context.Context, clientID, tagID int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("client_id = ? AND tag_id = ?", clientID, tagID).Delete(&entities.ClientClientTag{}).Error
}
