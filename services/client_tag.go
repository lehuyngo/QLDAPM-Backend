package services

import (
	"context"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClientTag interface {
	Create(ctx context.Context, data *entities.ClientTag) (int64, error)
	Add(ctx context.Context, data *entities.ClientClientTag) error
	Read(ctx context.Context, id int64) (*entities.ClientTag, error)
	ReadByUUID(ctx context.Context, orgID int64, uuid string) (*entities.ClientTag, error)
	ReadByName(ctx context.Context, orgID int64, name string) (*entities.ClientTag, error)
	Update(ctx context.Context, data *entities.ClientTag) error
	List(ctx context.Context, orgID int64) ([]*entities.ClientTag, error)
	Delete(ctx context.Context, user *entities.User, clientID, tagID int64) error
}

type ClientTag struct {
	Model         models.IClientTag
	JoinModel     models.IClientClientTag
	ClientModel   models.IClient
	ActivityModel models.IClientTagActivity
}

func NewClientTag() IClientTag {
	return &ClientTag{
		Model:         models.ClientTag{},
		JoinModel:     models.ClientClientTag{},
		ClientModel:   models.Client{},
		ActivityModel: models.ClientTagActivity{},
	}
}

func (p *ClientTag) Create(ctx context.Context, data *entities.ClientTag) (int64, error) {
	data.UUID = uuid.NewString()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	var clientID int64
	if len(data.Clients) > 0 {
		clientID = data.Clients[0].ClientID
		p.ClientModel.UpdateLastActiveTime(ctx, data.Clients[0].ClientID)
	}

	// add client tags activity
	p.ActivityModel.Create(ctx, &entities.ClientTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		TagID:     data.ID,
		Type:      entities.ActivityCreated,
		ClientID:  clientID,
	})

	return result, nil
}

func (p *ClientTag) Add(ctx context.Context, data *entities.ClientClientTag) error {
	err := p.JoinModel.Create(ctx, data)
	if err != nil {
		return err
	}

	p.ClientModel.UpdateLastActiveTime(ctx, data.ClientID)

	// add client tags activity
	p.ActivityModel.Create(ctx, &entities.ClientTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		TagID:     data.TagID,
		Type:      entities.ActivityCreated,
		ClientID:  data.ClientID,
	})

	return nil
}

func (p *ClientTag) Read(ctx context.Context, id int64) (*entities.ClientTag, error) {
	return p.Model.First(ctx, map[string]any{"id": id})
}

func (p *ClientTag) ReadByUUID(ctx context.Context, orgID int64, uuid string) (*entities.ClientTag, error) {
	return p.Model.First(ctx, map[string]any{"organization_id": orgID, "uuid": uuid})
}

func (p *ClientTag) ReadByName(ctx context.Context, orgID int64, name string) (*entities.ClientTag, error) {
	return p.Model.First(ctx, map[string]any{"organization_id": orgID, "name": name})
}

func (p *ClientTag) Update(ctx context.Context, data *entities.ClientTag) error {
	return p.Model.Update(ctx, data)
}

func (p *ClientTag) List(ctx context.Context, orgID int64) ([]*entities.ClientTag, error) {
	return p.Model.List(ctx, map[string]any{"organization_id": orgID})
}

func (p *ClientTag) Delete(ctx context.Context, user *entities.User, clientID, tagID int64) error {

	p.Model.Delete(ctx, tagID)

	p.ClientModel.UpdateLastActiveTime(ctx, clientID)

	// add client tags activity
	p.ActivityModel.Create(ctx, &entities.ClientTagActivity{
		UUID:      uuid.NewString(),
		CreatedBy: user.ID,
		TagID:     tagID,
		Type:      entities.ActivityDeleted,
		ClientID:  clientID,
	})

	return nil
}
