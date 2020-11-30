package lib

type causer interface {
	Cause() error
}

// ErrorCauser returns the error causer, works with wrapped errors.
func ErrorCauser(e error) error {
	causer, ok := e.(causer)
	if ok {
		return causer.Cause()
	}
	return e
}

// ErrorAPI contains an error_message, returned by server.
type ErrorAPI struct {
	ErrorMessage string
}

// NewErrorAPI ...
func NewErrorAPI(errMessage string) *ErrorAPI {
	return &ErrorAPI{
		ErrorMessage: errMessage,
	}
}

// Error ...
func (e *ErrorAPI) Error() string {
	return e.ErrorMessage
}

// IsErrorAPI ...
func IsErrorAPI(e error) bool {
	_, ok := ErrorCauser(e).(*ErrorAPI)
	return ok
}

// ------------------------------------------------------------------------

// ErrorInvalidEnum denotes that a enum/const variable holds invalid value.
type ErrorInvalidEnum struct{}

// NewErrorInvalidEnum ...
func NewErrorInvalidEnum() *ErrorInvalidEnum {
	return &ErrorInvalidEnum{}
}

// Error ...
func (ErrorInvalidEnum) Error() string {
	return "invalid enum value"
}

// IsErrorInvalidEnum ...
func IsErrorInvalidEnum(e error) bool {
	_, ok := ErrorCauser(e).(*ErrorInvalidEnum)
	return ok
}

// -------------------------------------------------------------------------

// ErrorInvalidArgument contains the argument that has caused the error.
type ErrorInvalidArgument struct {
	Argument string
}

// NewErrorInvalidArgument ...
func NewErrorInvalidArgument(arg string) *ErrorInvalidArgument {
	return &ErrorInvalidArgument{
		Argument: arg,
	}
}

// Error ...
func (e *ErrorInvalidArgument) Error() string {
	return e.Argument
}

// IsErrorInvalidArgument ...
func IsErrorInvalidArgument(e error) bool {
	_, ok := ErrorCauser(e).(*ErrorInvalidArgument)
	return ok
}
