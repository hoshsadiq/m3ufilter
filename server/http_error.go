package server

type HttpError interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (s *StatusError) Error() string {
	return s.Err.Error()
}

func (s *StatusError) Status() int {
	return s.Code
}
