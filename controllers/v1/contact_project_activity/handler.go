package contact_project_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                   services.IUser
	Contact                services.IContact
	Project                services.IProject
	ContactProjectActivity services.IContactProjectActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                   services.NewUser(),
		Contact:                services.NewContact(),
		Project:                services.NewProject(),
		ContactProjectActivity: services.NewContactProjectActivity(),
	}
}
