package receive_mail_history

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	TrackedURL services.ITrackedURL
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		TrackedURL: services.NewTrackedURL(),
	}
}
