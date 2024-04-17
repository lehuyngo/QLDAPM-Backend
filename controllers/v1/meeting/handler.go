package meeting

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User        services.IUser
	Meeting     services.IMeeting
	MeetingNote services.IMeetingNote
	Project     services.IProject

	Attendee    services.IAttendee
}

func NewHandler() *Handler {
	return &Handler{
		User:        services.NewUser(),
		Meeting:     services.NewMeeting(),
		MeetingNote: services.NewMeetingNote(),
		Project:     services.NewProject(),

		Attendee:    services.NewAttendee(),
	}
}
