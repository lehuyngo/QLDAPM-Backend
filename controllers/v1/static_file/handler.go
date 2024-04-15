package static_file

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User services.IUser
	StaticFile services.IStaticFile
}

func NewHandler() *Handler {
	return &Handler{
		User: services.NewUser(),
		StaticFile: services.NewStaticFile(),
	}
}
