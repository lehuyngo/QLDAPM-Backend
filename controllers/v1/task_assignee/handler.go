package task_assignee

import "gitlab.tgl-cloud.com/dx-ecosystem/crm/services"

type Handler struct {
	User         services.IUser
	Task         services.ITask
	TaskAssignee services.ITaskAssignee
	Project      services.IProject
}

func NewHandler() *Handler {
	return &Handler{
		User:         services.NewUser(),
		Task:         services.NewTask(),
		TaskAssignee: services.NewTaskAssignee(),
		Project:      services.NewProject(),
	}
}
