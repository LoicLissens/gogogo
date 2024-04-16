package customerrors

import "errors"

type ErrorNotFound struct {
	error
}

func NewErrorNotFound(message string) ErrorNotFound {
	return ErrorNotFound{errors.New(message)}
}

type ErrorAlreadyExists struct {
	error
}

func NewErrorAlreadyExists(message string) ErrorAlreadyExists {
	return ErrorAlreadyExists{errors.New(message)}
}
