package contact_note

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	ContactNote services.IContactNote
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		ContactNote: services.NewContactNote(),
	}
}
