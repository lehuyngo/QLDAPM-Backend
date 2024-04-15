package contact_project

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	Project services.IProject
	ContactProject services.IContactProject
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		Project: services.NewProject(),
		ContactProject: services.NewContactProject(),
	}
}
