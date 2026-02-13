package errs

import "fmt"

const InternalType = "INTERNAL"
const ValidationType = "VALIDATION"
const BadRequestType = "BAD_REQUEST"
const NotFoundType = "NOT_FOUND"
const UnauthorizedType = "UNAUTHORIZED"
const ForbiddenType = "FORBIDDEN"

type Error struct {
	Code    string
	Field   string
	Message string
	Type    string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(code, msg, typ string, err error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Type:    typ,
		Err:     err,
	}
}

func NewFieldError(code, field, msg, typ string, err error) *Error {
	return &Error{
		Code:    code,
		Field:   field,
		Message: msg,
		Type:    typ,
		Err:     err,
	}
}

type ErrorList struct {
	Errors []*Error
}

func (l *ErrorList) Error() string {
	if len(l.Errors) == 0 {
		return "no errors"
	}
	if len(l.Errors) == 1 {
		return l.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred", len(l.Errors))
}

func (l *ErrorList) Append(err *Error) {
	if err == nil {
		return
	}
	l.Errors = append(l.Errors, err)
}

func (l *ErrorList) Len() int {
	return len(l.Errors)
}

func NewErrorList(errs ...*Error) *ErrorList {
	list := &ErrorList{}
	for _, err := range errs {
		list.Append(err)
	}
	return list
}
