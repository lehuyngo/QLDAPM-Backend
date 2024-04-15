package meeting_highlight

import (
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

type Handler struct {
	User        services.IUser
	MeetingNote services.IMeetingNote
	Highlight   services.IMeetingHighlight
	Project     services.IProject
}

func NewHandler() *Handler {
	return &Handler{
		User:        services.NewUser(),
		MeetingNote: services.NewMeetingNote(),
		Highlight:   services.NewMeetingHighlight(),
		Project:     services.NewProject(),
	}
}
