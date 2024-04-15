package models

import (
	"context"
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/clients"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

type IClientContact interface {
	Create(ctx context.Context, data *entities.ClientContact) error
	CreateBatch(ctx context.Context, data []entities.ClientContact) error
	Delete(ctx context.Context, clientID, contactID int64) error
	List(ctx context.Context, filters map[string]any) ([]*entities.ClientContact, error)
}

type ClientContact struct {
}

func (ClientContact) Create(ctx context.Context, data *entities.ClientContact) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(data).Error
}

func (ClientContact) CreateBatch(ctx context.Context, data []entities.ClientContact) error {
	db := clients.MySQLClient
	return db.WithContext(ctx).Create(data).Error
}

func (ClientContact) Delete(ctx context.Context, clientID, contactID int64) error {
	db := clients.MySQLClient.WithContext(ctx)
	return db.Where("client_id = ? AND contact_id = ?", clientID, contactID).Delete(&entities.ClientContact{}).Error
}

func (ClientContact) List(ctx context.Context, filters map[string]any) ([]*entities.ClientContact, error) {
	if len(filters) < 1 {
		return nil, define.ErrFilterIsEmpty
	}
	
	result := []*entities.ClientContact{}

	db := clients.MySQLClient.WithContext(ctx)
	for field, val := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", field), val)
	}

	err := db.Preload("Client.Logo").
		Preload("Contact.Avatar").
		Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
