package task_attach_file

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User           services.IUser
	Task           services.ITask
	TaskAttachFile services.ITaskAttachFile
}

func NewHandler() *Handler {
	return &Handler{
		User:           services.NewUser(),
		Task:           services.NewTask(),
		TaskAttachFile: services.NewTaskAttachFile(),
	}
}
