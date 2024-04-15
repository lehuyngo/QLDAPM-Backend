package contact_tag_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User               services.IUser
	Contact            services.IContact
	ContactTagActivity services.IContactTagActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:               services.NewUser(),
		Contact:            services.NewContact(),
		ContactTagActivity: services.NewContactTagActivity(),
	}
}
