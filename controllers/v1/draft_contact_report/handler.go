package draft_contact_report

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	DraftContactReport services.IDraftContactReport
}
func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		DraftContactReport: services.NewDraftContactReport(),
	}
}