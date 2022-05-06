package logger

import "fmt"

type Logger interface {
	Info(string)
	Error(string)
	Debug(string)
	Warn(string)
}

type DefaultLogger struct{}

func (log *DefaultLogger) Info(msg string) {
	fmt.Println(msg)
}

func (log *DefaultLogger) Error(msg string) {
	fmt.Println(msg)
}

func (log *DefaultLogger) Debug(msg string) {
	fmt.Println(msg)
}

func (log *DefaultLogger) Warn(msg string) {
	fmt.Println(msg)
}
