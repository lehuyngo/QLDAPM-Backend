package client

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Client services.IClient
	ClientContact services.IClientContact
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Client: services.NewClient(),
		ClientContact: services.NewClientContact(),
	}
}
