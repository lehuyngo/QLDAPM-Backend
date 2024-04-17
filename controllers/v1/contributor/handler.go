package contributor

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User        services.IUser
	
	Contributor services.IContributor
	MeetingNote services.IMeetingNote
}

func NewHandler() *Handler {
	return &Handler{
		User:        services.NewUser(),

		Contributor: services.NewContributor(),
		MeetingNote: services.NewMeetingNote(),
	}
}
