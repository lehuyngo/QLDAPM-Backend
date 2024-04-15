package contributor

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User        services.IUser
	Contact     services.IContact
	Contributor services.IContributor
	MeetingNote services.IMeetingNote
}

func NewHandler() *Handler {
	return &Handler{
		User:        services.NewUser(),
		Contact:     services.NewContact(),
		Contributor: services.NewContributor(),
		MeetingNote: services.NewMeetingNote(),
	}
}
