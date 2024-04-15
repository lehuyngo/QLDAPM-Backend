package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/util"
)

type IUser interface {
	Register(ctx context.Context, data *entities.User) (string, error)
	Create(ctx context.Context, data *entities.User) (string, error)
	Read(ctx context.Context, id int64) (*entities.User, error)
	ReadByUUID(ctx context.Context, uuid string) (*entities.User, error)
	ReadByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, data *entities.User) error
	Delete(ctx context.Context, id int64) error
	Authenticate(ctx context.Context, username, password string) (string, error)
	ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.User, error)
	ListByOrgID(ctx context.Context, orgID int64) ([]*entities.User, error)
}

type User struct {
	Model models.IUser
}

func NewUser() IUser {
	return &User{
		Model: models.User{},
	}
}

func (p *User) Register(ctx context.Context, data *entities.User) (string, error) {
	data.PasswordHash = util.GetMD5Hash(data.PasswordHash)
	data.UUID = uuid.NewString()
	// data.Organization.UUID = uuid.NewString()
	userId, err := p.Model.Create(ctx, data)
	if err != nil {
		return "", err
	}

	data.Token, err = GenerateToken(data.ID, time.Now().UnixNano())
	if err != nil {
		p.Model.Delete(ctx, data.ID)
		return "", err
	}

	data.OrganizationID = data.Organization.ID
	data.CreatedBy = data.ID
	data.UpdatedBy = data.ID
	data.Organization.CreatedBy = data.ID
	data.Organization.UpdatedBy = data.ID

	err = p.Model.Update(ctx, data)
	if err != nil {
		p.Model.Delete(ctx, data.ID)
		return "", err
	}

	return TokenMaker.CreateToken(userId)
}

func (p *User) Create(ctx context.Context, data *entities.User) (string, error) {
	data.PasswordHash = util.GetMD5Hash(data.PasswordHash)
	data.UUID = uuid.NewString()
	data.Status = entities.Active.Value()
	data.Organization.UUID = uuid.NewString()
	userId, err := p.Model.Create(ctx, data)
	if err != nil {
		return "", err
	}

	data.Token, err = GenerateToken(data.ID, time.Now().UnixNano())
	if err != nil {
		p.Model.Delete(ctx, data.ID)
		return "", err
	}

	err = p.Model.Update(ctx, data)
	if err != nil {
		p.Model.Delete(ctx, data.ID)
		return "", err
	}

	return TokenMaker.CreateToken(userId)
}

func (p *User) Read(ctx context.Context, id int64) (*entities.User, error) {
	return p.Model.Read(ctx, id)
}

func (p *User) ReadByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	return p.Model.ReadByUUID(ctx, uuid)
}

func (p *User) ReadByUsername(ctx context.Context, username string) (*entities.User, error) {
	return p.Model.ReadByCondition(ctx, "username", username)
}

func (p *User) Update(ctx context.Context, data *entities.User) error {
	return p.Model.Update(ctx, data)
}

func (p *User) Delete(ctx context.Context, id int64) error {
	return p.Model.Delete(ctx, id)
}

func (p *User) Authenticate(ctx context.Context, username, password string) (string, error) {
	data, err := p.Model.ReadByCondition(ctx, "username", username)
	if err != nil {
		return "", err
	}

	if util.GetMD5Hash(password) != data.PasswordHash {
		return "", fmt.Errorf("Password is incorrect!")
	}

	return TokenMaker.CreateToken(data.ID)
}

func (p *User) ListByUUIDs(ctx context.Context, uuids []string) ([]*entities.User, error) {
	return p.Model.ListByUUIDs(ctx, uuids)
}

func (p *User) ListByOrgID(ctx context.Context, orgID int64) ([]*entities.User, error) {
	return p.Model.ListByOrgID(ctx, orgID)
}
