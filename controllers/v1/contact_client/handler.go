package contact_client

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Client services.IClient
	Contact services.IContact
	ContactClient services.IContactClient
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Client: services.NewClient(),
		Contact: services.NewContact(),
		ContactClient: services.NewContactClient(),
	}
}
