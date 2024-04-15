package account

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Organization services.IOrganization
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Organization: services.NewOrganization(),
	}
}
