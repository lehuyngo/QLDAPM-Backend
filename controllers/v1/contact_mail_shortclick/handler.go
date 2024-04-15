package contact_mail_shortclick

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User            		services.IUser
	Contact         		services.IContact
	ContactMailShortClick	services.IContactMailShortClick
}

func NewHandler() *Handler {
	return &Handler{
		User:            		services.NewUser(),
		Contact:         		services.NewContact(),
		ContactMailShortClick: 	services.NewContactMailShortClick(),
	}
}
