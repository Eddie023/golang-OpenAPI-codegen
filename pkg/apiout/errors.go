package apiout

import "errors"

type BadRequestErr struct {
	Msg string
}

func (e *BadRequestErr) Error() string {
	return e.Msg
}

func BadRequest(msg string) error {
	return &BadRequestErr{Msg: msg}
}

func IsBadRequest(err error) bool {
	var be *BadRequestErr

	return errors.As(err, &be)
}
