package apis

type ErrorMessage struct {
	VI string `json:"vi,omitempty"`
	EN string `json:"en,omitempty"`
	JP string `json:"jp,omitempty"`
}

type Error struct {
	Domain  string        `json:"domain,omitempty"`
	Reason  string        `json:"reason,omitempty"`
	Message *ErrorMessage `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
}

type Empty struct {
}

type CreateResponse struct {
	UUID string `json:"uuid"`
}

type TimeRange struct {
	StartTime int64 `json:"start_time,omitempty"`
	EndTime   int64 `json:"end_time,omitempty"`
}

type File struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Paging struct {
	Page int `form:"page,omitempty"`
	Size int `form:"size,omitempty"`
}

type ListFileResponse struct {
	Data []*File `json:"data"`
}

type ActivityType int

const (
	ActivityCreated ActivityType = 0
	ActivityDeleted ActivityType = 1
	ActivityUpdated ActivityType = 2
)

type TaskStatus int

const (
	TaskStatusToDo    TaskStatus = 1
	TaskStatusDoing   TaskStatus = 2
	TaskStatusTesting TaskStatus = 3
	TaskStatusDone    TaskStatus = 4
)

type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 1
	TaskPriorityMedium TaskPriority = 2
	TaskPriorityHigh   TaskPriority = 3
)

type TaskLabel int

const (
	TaskLabelTask        TaskLabel = 1
	TaskLabelFeedback    TaskLabel = 2
	TaskLabelImprovement TaskLabel = 3
	TaskLabelBug         TaskLabel = 4
)
