package contact_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User            services.IUser
	Contact         services.IContact
	ContactActivity services.IContactActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:            services.NewUser(),
		Contact:         services.NewContact(),
		ContactActivity: services.NewContactActivity(),
	}
}
