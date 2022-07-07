package service

type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
}

type nopLogger struct{}

func (nopLogger) Error(...interface{})          {}
func (nopLogger) Errorf(string, ...interface{}) {}
