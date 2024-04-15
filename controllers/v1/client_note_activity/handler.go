package client_note_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User               services.IUser
	Client             services.IClient
	ClientNoteActivity services.IClientNoteActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:               services.NewUser(),
		Client:             services.NewClient(),
		ClientNoteActivity: services.NewClientNoteActivity(),
	}
}
