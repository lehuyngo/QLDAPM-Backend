package entities

type ProjectStatus int32

const (
	Prospect		ProjectStatus = 1
	Contacted		ProjectStatus = 2
	Estimate		ProjectStatus = 3
	FollowUp		ProjectStatus = 4
	ProjectReceived	ProjectStatus = 5
)

func (s ProjectStatus) Value() int32 {
	return int32(s)
}

var mapProjectStatusName = map[ProjectStatus]string {
	Prospect: "Prospect",
	Contacted: "Contacted",
	Estimate: "Estimate",
	FollowUp: "Follow up",
	ProjectReceived: "Project received",
}

var mapProjectStatusNote = map[ProjectStatus]string {
	ProjectReceived: "If user select this status, the project will be moved to PIMS page. You can not change the status again.",
}

func (s ProjectStatus) Name() string {
	val, exits := mapProjectStatusName[s]
	if !exits {
		return ""
	}

	return val
}

func (s ProjectStatus) Note() string {
	val, exits := mapProjectStatusNote[s]
	if !exits {
		return ""
	}

	return val
}