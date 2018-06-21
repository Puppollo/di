package di

import (
	"fmt"
	"os"
)

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
}

type (
	nullLogger    struct{}
	consoleLogger struct{}
)

func (consoleLogger) Debug(v ...interface{}) {
	fmt.Println(v...)
}

func (consoleLogger) Info(v ...interface{}) {
	fmt.Println(v...)
}

func (consoleLogger) Error(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func (nullLogger) Debug(...interface{}) {}
func (nullLogger) Info(...interface{})  {}
func (nullLogger) Error(...interface{}) {}

func NullLogger() Logger {
	return nullLogger{}
}

func ConsoleLogger() Logger {
	return consoleLogger{}
}
