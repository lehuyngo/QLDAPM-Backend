package apis

type Task struct {
	UUID           string       `json:"uuid,omitempty"`
	Title          string       `json:"title,omitempty"`
	Status         TaskStatus   `json:"status,omitempty"`
	Priority       TaskPriority `json:"priority,omitempty"`
	Label          TaskLabel    `json:"label,omitempty"`
	ProjectID      int64        `json:"project_id,omitempty"`
	Project        *Project     `json:"project,omitempty"`
	Deadline       int64        `json:"deadline,omitempty"`
	EstimatedHours int32        `json:"estimated_hours,omitempty"`
	Description    string       `json:"description,omitempty"`
	Assignees      []*User      `json:"assignees,omitempty"`
	CreatedTime    int64        `json:"created_time,omitempty"`
	Creator        *User        `json:"creator,omitempty"`
	AttachFiles    []*File      `json:"attach_files,omitempty"`
}

func (c *Task) GetUUID() string {
	if c == nil {
		return ""
	}

	return c.UUID
}

type ListTaskResponse struct {
	Data []*Task `json:"data"`
}

type CreateTaskRequest struct {
	Title          string       `form:"title,omitempty"`
	Status         TaskStatus   `form:"status,omitempty"`
	Priority       TaskPriority `form:"priority,omitempty"`
	Label          TaskLabel    `form:"label,omitempty"`
	ProjectUUID    string       `form:"project_uuid,omitempty"`
	Deadline       int64        `form:"deadline,omitempty"`
	EstimatedHours int32        `form:"estimated_hours,omitempty"`
	Description    string       `form:"description,omitempty"`
	AssigneeUUIDs  string       `form:"assignee_uuids,omitempty"`
}

type UpdateTaskRequest struct {
	Title          string       `json:"title,omitempty"`
	Status         TaskStatus   `json:"status,omitempty"`
	Priority       TaskPriority `json:"priority,omitempty"`
	Label          TaskLabel    `json:"label,omitempty"`
	ProjectUUID    string       `json:"project_uuid,omitempty"`
	Deadline       int64        `json:"deadline,omitempty"`
	EstimatedHours int32        `json:"estimated_hours,omitempty"`
	Description    string       `json:"description,omitempty"`
}

type UpdateTaskStatusRequest struct {
	Status TaskStatus `json:"status,omitempty"`
}
