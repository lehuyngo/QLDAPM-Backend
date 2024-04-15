package task

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User           services.IUser
	Task           services.ITask
	TaskAssignee   services.ITaskAssignee
	Project        services.IProject
	TaskAttachFile services.ITaskAttachFile
}

func NewHandler() *Handler {
	return &Handler{
		User:           services.NewUser(),
		Task:           services.NewTask(),
		TaskAssignee:   services.NewTaskAssignee(),
		Project:        services.NewProject(),
		TaskAttachFile: services.NewTaskAttachFile(),
	}
}
