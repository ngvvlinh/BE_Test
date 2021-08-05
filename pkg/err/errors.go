package err

import (
	"fmt"
)

func New(num int, text string) *ErrorString {
	return &ErrorString{num, text}
}

type ErrorString struct {
	n int
	s string
}

func (e *ErrorString) Error() string {
	return e.s
}

func (e *ErrorString) String() string {
	return e.s
}

func (e *ErrorString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"%d\":\"%s\"}", e.n, e.s)), nil
}

var (
	ErrSuccess          = New(0, "Success")
	ErrInvalid          = New(1, "Invalid function")
	ErrListInconsistent = New(255, "Extended attributes are inconsistent")

	Errors = map[int]*ErrorString{
		0:   ErrSuccess,
		1:   ErrInvalid,
		255: ErrListInconsistent,
	}
)

func ErrorsList() []*ErrorString {
	errs := make([]*ErrorString, len(Errors))

	i := 0
	for idx := range Errors {
		errs[i] = Errors[idx]
		i++
	}

	return errs
}
