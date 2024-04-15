package contact_client_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                  services.IUser
	Contact               services.IContact
	Client                services.IClient
	ContactClientActivity services.IContactClientActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                  services.NewUser(),
		Contact:               services.NewContact(),
		Client:                services.NewClient(),
		ContactClientActivity: services.NewContactClientActivity(),
	}
}
