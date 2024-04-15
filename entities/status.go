package entities

type Status int32

const (
	InTrash		Status = -1
	Inactive	Status = 0
	Active 		Status = 1
)

func (s Status) Value() int32 {
	return int32(s)
}