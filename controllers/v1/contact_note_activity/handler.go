package contact_note_activity

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User                services.IUser
	Contact             services.IContact
	ContactNoteActivity services.IContactNoteActivity
}

func NewHandler() *Handler {
	return &Handler{
		User:                services.NewUser(),
		Contact:             services.NewContact(),
		ContactNoteActivity: services.NewContactNoteActivity(),
	}
}
