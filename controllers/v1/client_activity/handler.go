package client_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User           services.IUser
	Client         services.IClient
	ClientActivity services.IClientActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:           services.NewUser(),
		Client:         services.NewClient(),
		ClientActivity: services.NewClientActivity(),
	}
}
