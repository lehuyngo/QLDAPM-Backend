package report_read_mail

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	EventClickURL services.IEventClickURL
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		EventClickURL: services.NewEventClickURL(),
	}
}
