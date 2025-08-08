package error

import "fmt"

type DuplicateError struct {
	Msg string
}

func NewDuplicateError(msg string) *DuplicateError {
	return &DuplicateError{
		Msg: msg,
	}
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf(e.Msg)
}
