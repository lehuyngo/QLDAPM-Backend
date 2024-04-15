package services

import (
	"context"
	"errors"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gorm.io/gorm"
)

type IClientClientTag interface {
	Delete(ctx context.Context, clientID, tagID int64) error
}

type ClientClientTag struct {
	Model          models.IClientClientTag
	ClientTagModel models.IClientTag
}

func NewClientClientTag() IClientClientTag {
	return &ClientClientTag{
		Model:          models.ClientClientTag{},
		ClientTagModel: models.ClientTag{},
	}
}

func (p *ClientClientTag) Delete(ctx context.Context, clientID, tagID int64) error {

	err := p.Model.Delete(ctx, clientID, tagID)
	if err != nil {
		return err
	}
	_, err = p.Model.First(ctx, tagID)

	// remove untagged tag
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return p.ClientTagModel.Delete(ctx, tagID)
	}

	return nil
}
