package draft_contact

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	Contact services.IContact
	Client services.IClient
	ClientTag services.IClientTag
	DraftContact services.IDraftContact
	ContactClient services.IContactClient
	ContactTag services.IContactTag
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		Contact: services.NewContact(),
		Client: services.NewClient(),
		ClientTag: services.NewClientTag(),
		DraftContact: services.NewDraftContact(),
		ContactClient: services.NewContactClient(),
		ContactTag: services.NewContactTag(),
	}
}
