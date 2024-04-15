package apis

type TaskAssignee struct {
	TaskID      int64 `json:"task_id,omitempty"`
	Task        *Task `json:"task,omitempty"`
	UserID      int64 `json:"user_id,omitempty"`
	User        *User `json:"user,omitempty"`
	CreatedTime int64 `json:"created_time,omitempty"`
	Creator     *User `json:"creator,omitempty"`
}

func (c *TaskAssignee) GetUser() *User {
	if c == nil {
		return nil
	}

	return c.User
}

type CreateTaskAssigneeRequest struct {
	AssigneeUUID string `json:"assignee_uuid,omitempty"`
}
