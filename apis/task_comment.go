package apis

type TaskComment struct {
	UUID        string `json:"uuid,omitempty"`
	Content     string `json:"content,omitempty"`
	Creator     *User  `json:"creator,omitempty"`
	CreatedTime int64  `json:"created_time,omitempty"`
}
type CreateTaskCommentRequest struct {
	Content string `json:"content,omitempty"`
}
type ListTaskComment struct {
	Data []*TaskComment `json:"data,omitempty"`
}