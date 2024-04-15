package meeting_note

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User              services.IUser
	Meeting           services.IMeeting
	MeetingNote       services.IMeetingNote
	Project           services.IProject
	Contact           services.IContact
	Contributor       services.IContributor
	MeetingNoteEditor services.IMeetingNoteEditor
}

func NewHandler() *Handler {
	return &Handler{
		User:              services.NewUser(),
		Meeting:           services.NewMeeting(),
		MeetingNote:       services.NewMeetingNote(),
		Project:           services.NewProject(),
		Contact:           services.NewContact(),
		Contributor:       services.NewContributor(),
		MeetingNoteEditor: services.NewMeetingNoteEditor(),
	}
}
