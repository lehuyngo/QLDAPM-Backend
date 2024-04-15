package entities

type Gender int32

const (
	Female	Gender = 1
	Male	Gender = 2
	Other	Gender = 3
)

func (s Gender) Value() int32 {
	return int32(s)
}
