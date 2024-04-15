package contact_mail_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                services.IUser
	Contact             services.IContact
	ContactMailActivity services.IContactMailActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                services.NewUser(),
		Contact:             services.NewContact(),
		ContactMailActivity: services.NewContactMailActivity(),
	}
}
