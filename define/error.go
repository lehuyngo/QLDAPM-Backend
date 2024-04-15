package define

import (
	"encoding/json"
	"errors"
)

var ErrFilterIsEmpty = errors.New("filter is empty")
var ErrNotInitial = errors.New("not initial")

type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

func IsErrDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	
	bytes, e := json.Marshal(err)
	if e != nil {
		return false
	}

	var newError GormErr
	e = json.Unmarshal((bytes), &newError)
	if e != nil {
		return false
	}

	return (newError.Number == 1062)
}