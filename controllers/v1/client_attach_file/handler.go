package client_attach_file

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Client services.IClient
	ClientAttachFile services.IClientAttachFile
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Client: services.NewClient(),
		ClientAttachFile: services.NewClientAttachFile(),
	}
}
