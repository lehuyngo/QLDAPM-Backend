package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
)

type IClient interface {
	Create(ctx context.Context, data *entities.Client) (int64, error)
	Read(ctx context.Context, id int64) (*entities.Client, error)
	ReadByName(ctx context.Context, orgID int64, name string) (*entities.Client, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.Client, error)
	ReadByWebsite(ctx context.Context, orgID int64, website string) (*entities.Client, error)
	Update(ctx context.Context, data *entities.Client) error
	Delete(ctx context.Context, userID, id int64) error
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Client, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Client, error)
}

type Client struct {
	Model         models.IClient
	ActivityModel models.ClientActivity
}

func NewClient() IClient {
	return &Client{
		Model:         models.Client{},
		ActivityModel: models.ClientActivity{},
	}
}

func (p *Client) Create(ctx context.Context, data *entities.Client) (int64, error) {
	data.UUID = uuid.NewString()
	data.Status = entities.Active
	data.LastActiveTime = time.Now().UnixMilli()
	result, err := p.Model.Create(ctx, data)
	if err != nil {
		return 0, err
	}
	p.ActivityModel.Create(ctx, &entities.ClientActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		ClientID:  data.ID,
		Type:      entities.ActivityCreated,
	})
	return result, nil
}

func (p *Client) Read(ctx context.Context, id int64) (*entities.Client, error) {
	return p.Model.First(ctx, map[string]any{"id": id})
}

func (p *Client) ReadByName(ctx context.Context, orgID int64, name string) (*entities.Client, error) {
	filters := map[string]any{"organization_id": orgID, "fullname": name}
	return p.Model.First(ctx, filters)
}

func (p *Client) ReadByUUID(ctx context.Context, uuid string) (*entities.Client, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *Client) ReadByWebsite(ctx context.Context, orgID int64, website string) (*entities.Client, error) {
	filters := map[string]any{"website": website, "organization_id": orgID}
	return p.Model.First(ctx, filters)
}

func (p *Client) Update(ctx context.Context, data *entities.Client) error {
	data.LastActiveTime = time.Now().UnixMilli()
	err := p.Model.Update(ctx, data)
	if err != nil {
		return err
	}

	typeActivity := entities.ActivityUpdated
	if data.Status == entities.InTrash {
		typeActivity = entities.ActivityDeleted
	}
	p.ActivityModel.Create(ctx, &entities.ClientActivity{
		UUID:      uuid.NewString(),
		CreatedBy: data.CreatedBy,
		ClientID:  data.ID,
		Type:      typeActivity,
	})
	return nil
}

func (p *Client) Delete(ctx context.Context, userID, id int64) error {
	err := p.Model.Delete(ctx, id)
	if err != nil {
		return err
	}
	p.ActivityModel.Create(ctx, &entities.ClientActivity{
		UUID:      uuid.NewString(),
		ClientID:  id,
		CreatedBy: userID,
		Type:      entities.ActivityDeleted,
	})
	return nil
}

func (p *Client) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.Client, error) {
	filters := map[string]any{"organization_id": orgID}
	return p.Model.List(ctx, filters)
}

func (p *Client) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.Client, error) {
	return p.Model.ListByUUIDs(ctx, uuids)
}
