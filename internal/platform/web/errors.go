package web

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type Error struct {
	Err error
	Status int
	Fields []FieldError
}

func NewRequestError(err error, status int) error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
