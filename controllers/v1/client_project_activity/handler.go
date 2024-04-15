package client_project_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                  services.IUser
	Client                services.IClient
	Project               services.IProject
	ClientProjectActivity services.IClientProjectActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                  services.NewUser(),
		Client:                services.NewClient(),
		Project:               services.NewProject(),
		ClientProjectActivity: services.NewClientProjectActivity(),
	}
}
