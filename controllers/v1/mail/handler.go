package mail

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	Mail services.IMail
	TrackedURL services.ITrackedURL
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		Mail: services.NewMail(),
		TrackedURL: services.NewTrackedURL(),
	}
}
