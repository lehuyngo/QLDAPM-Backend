package client_contact_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                  services.IUser
	Contact               services.IContact
	Client                services.IClient
	ClientContactActivity services.IClientContactActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                  services.NewUser(),
		Contact:               services.NewContact(),
		Client:                services.NewClient(),
		ClientContactActivity: services.NewClientContactActivity(),
	}
}
