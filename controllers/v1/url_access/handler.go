package url_access

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	TrackedURL services.ITrackedURL
	EventClickURL services.IEventClickURL
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		TrackedURL: services.NewTrackedURL(),
		EventClickURL: services.NewEventClickURL(),
	}
}
