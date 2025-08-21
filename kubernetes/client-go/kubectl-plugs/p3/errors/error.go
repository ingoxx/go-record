package errors

import (
	pe "errors"
	"fmt"
)

var (
	MissDataError = pe.New("miss data")
	NotFoundError = pe.New("namespace not found")
)

type DeleteError struct {
	Msg error
}

func (e DeleteError) Error() string {
	return e.Msg.Error()
}

func NewDeleteError(name string) error {
	return DeleteError{
		Msg: fmt.Errorf("fail to delete namespace '%s', maybe not marked for deletion", name),
	}
}
