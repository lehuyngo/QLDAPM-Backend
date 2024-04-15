package contact_tag

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	ContactNote services.IContactNote
	ContactTag services.IContactTag
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		ContactNote: services.NewContactNote(),
		ContactTag: services.NewContactTag(),
	}
}
