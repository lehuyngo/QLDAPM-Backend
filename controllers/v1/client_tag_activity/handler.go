package client_tag_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Client services.IClient
	ClientTagActivity services.IClientTagActivity
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Client: services.NewClient(),
		ClientTagActivity: services.NewClientTagActivity(),
	}
}
