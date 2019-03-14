package loki

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
	"time"
)

var (
	logger = &Logger{
		level:     INFO,
		formatter: NewStandardFormatter(),
		handler:   NewConsoleHandler(),
	}

	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
)

func SetLevel(level int) {
	logger.level = level
}

func SetFormatter(formatter Formatter) {
	logger.formatter = formatter
}

type Logger struct {
	level     int
	formatter Formatter
	handler   Handler
}

func Debug(format string, a ...interface{}) {
	if DEBUG >= logger.level {
		logger.handler.output(logger.formatter.format(format, a...))
	}
}

func Info(format string, a ...interface{}) {
	if INFO >= logger.level {
		logger.handler.output(aurora.Blue(logger.formatter.format(format, a...)))
	}
}

func Warn(format string, a ...interface{}) {
	if WARN >= logger.level {
		logger.handler.output(aurora.Green(logger.formatter.format(format, a...)))
	}
}

func Error(format string, a ...interface{}) {
	if ERROR >= logger.level {
		logger.handler.output(aurora.Red(logger.formatter.format(format, a...)))
	}
}

func Fatal(format string, a ...interface{}) {
	Error(format, a...)
	os.Exit(1)
}

type Formatter interface {
	format(format string, a ...interface{}) string
}

type StandardFormatter struct {
	Formatter
}

func NewStandardFormatter() Formatter {
	return StandardFormatter{}
}

func (f StandardFormatter) format(format string, a ...interface{}) string {
	return fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
}

type Handler interface {
	output(output interface{}) error
}

type ConsoleHandler struct {
	Handler
}

func NewConsoleHandler() Handler {
	return ConsoleHandler{}
}

func (c ConsoleHandler) output(output interface{}) error {
	_, err := fmt.Println(output)
	return err
}
