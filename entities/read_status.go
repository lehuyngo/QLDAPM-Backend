package entities

type ReadStatus int32

const (
	New		ReadStatus = 1
	Read	ReadStatus = 2
)

func (s ReadStatus) Value() int32 {
	return int32(s)
}
