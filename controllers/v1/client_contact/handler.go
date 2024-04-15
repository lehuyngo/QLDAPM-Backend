package client_contact

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Client services.IClient
	Contact services.IContact
	ClientContact services.IClientContact
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Client: services.NewClient(),
		Contact: services.NewContact(),
		ClientContact: services.NewClientContact(),
	}
}
