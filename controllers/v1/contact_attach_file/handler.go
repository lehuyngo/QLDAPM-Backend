package contact_attach_file

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	ContactAttachFile services.IContactAttachFile
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		ContactAttachFile: services.NewContactAttachFile(),
	}
}
