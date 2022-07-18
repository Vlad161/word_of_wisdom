package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Sync()
}

type logger struct {
	log *zap.SugaredLogger
}

func New() *logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	return &logger{log: zapLogger.Sugar()}
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args)
}

func (l *logger) Fatal(args ...interface{}) {
	l.log.Fatal(args)
}

func (l *logger) Sync() {
	if err := l.log.Sync(); err != nil {
		log.Println(err)
	}
}
