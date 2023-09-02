package errors

type errors struct {
	Code int
	Msg  string
}

func (e *errors) Error() string {
	return e.Msg
}

func New(code int, err error) error {
	return &errors{
		Code: code,
		Msg:  err.Error(),
	}
}

const (
	UnAuthorized = 404
)
