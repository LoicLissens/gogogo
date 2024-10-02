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

type ValueError struct {
	error
}

func NewValueError(message string) ValueError {
	return ValueError{errors.New(message)}
}
