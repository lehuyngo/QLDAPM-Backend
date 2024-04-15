package project

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User           services.IUser
	Project        services.IProject
	Client         services.IClient
	Contact        services.IContact
	ContactProject services.IContactProject
	Meeting        services.IMeeting
}

func NewHandler() *Handler {
	return &Handler{
		User:           services.NewUser(),
		Project:        services.NewProject(),
		Client:         services.NewClient(),
		Contact:        services.NewContact(),
		ContactProject: services.NewContactProject(),
		Meeting:        services.NewMeeting(),
	}
}
