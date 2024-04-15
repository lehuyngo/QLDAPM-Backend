package client_tag

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User            services.IUser
	Client          services.IClient
	ClientNote      services.IClientNote
	ClientTag       services.IClientTag
	ClientClientTag services.IClientClientTag
}

func NewHandler() *Handler {
	return &Handler{
		User:            services.NewUser(),
		Client:          services.NewClient(),
		ClientNote:      services.NewClientNote(),
		ClientTag:       services.NewClientTag(),
		ClientClientTag: services.NewClientClientTag(),
	}
}
