package service

type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
}

type nopLogger struct{}

func (nopLogger) Error(...interface{})          {}
func (nopLogger) Errorf(string, ...interface{}) {}
