package client_contact_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                  services.IUser

	Client                services.IClient

}

func NewHandler() *Handler {
	return &Handler{
		User:                  services.NewUser(),
		Client:                services.NewClient(),
	
	}
}
