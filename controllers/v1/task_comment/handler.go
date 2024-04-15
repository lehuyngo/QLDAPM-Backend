package task_comment

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User        services.IUser
	TaskComment services.ITaskComment
	Task        services.ITask
}

func NewHandler() *Handler {
	return &Handler{
		User:        services.NewUser(),
		TaskComment: services.NewTaskComment(),
		Task:        services.NewTask(),
	}
}