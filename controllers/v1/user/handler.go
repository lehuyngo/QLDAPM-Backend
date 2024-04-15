package user

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
	}
}
