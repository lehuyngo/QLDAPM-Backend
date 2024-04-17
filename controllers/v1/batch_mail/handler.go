package batch_mail

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser

	BatchMail services.IBatchMail
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		
		BatchMail: services.NewBatchMail(),
	}
}
