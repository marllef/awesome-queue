package logger

import (
	"fmt"
	"log"
)

type LoggerInterface interface {
	Infof()
	Warnf()
	Errorf()
	Debugf()
}

type Logger struct {
	prefix string
	level  int
}

func Default() *Logger {
	return &Logger{
		level:  1,
		prefix: "",
	}
}

func (l *Logger) Infof(format string, values ...interface{}) {
	log.Printf(fmt.Sprintf("\033[36m[INFO] %s%s\033[0m", l.prefix, format), values...)
}

func (l *Logger) Warnf(format string, values ...interface{}) {
	log.Printf(fmt.Sprintf("\033[33m[WARN] %s%s\033[0m", l.prefix, format), values...)
}

func (l *Logger) Errorf(format string, values ...interface{}) {
	log.Printf(fmt.Sprintf("\033[31m[ERROR] %s%s\033[0m", l.prefix, format), values...)
}

func (l *Logger) Debugf(format string, values ...interface{}) {
	if l.level == 0 {
		log.Printf(fmt.Sprintf("[DEBUG] %s%s", l.prefix, format), values...)
	}
}
