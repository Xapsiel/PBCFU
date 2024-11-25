package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

type LogService struct {
	output io.Writer
	logger *logrus.Logger
}

var Logger *LogService

func NewLogService(output io.Writer) *LogService {
	// Создаем логгер с привязкой к контексту
	logger := logrus.New()
	logger.SetOutput(output)
	Logger = &LogService{
		output: output,
		logger: logger,
	}
	return Logger
}

// Print логирует сообщение с ID и контекстом
func (l *LogService) Print(id int, msg string) {
	// Получаем имя вызывающей функции
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	switch id {
	case -1:
		l.logger.WithFields(logrus.Fields{
			"fromFunc": callerName,
		}).Info(msg)
	default:
		l.logger.WithFields(logrus.Fields{
			"user_id":  id,
			"fromFunc": callerName,
		}).Info(msg)
	}

}
func (l *LogService) Warn(id int, msg string) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()
	switch id {
	case -1:
		l.logger.WithFields(logrus.Fields{
			"fromFunc": callerName,
		}).Warn(msg)
	default:
		l.logger.WithFields(logrus.Fields{
			"user_id":  id,
			"fromFunc": callerName,
		}).Warn(msg)
	}
}

func (l *LogService) SetFormat(writer io.Writer) {
	Logger.output = writer
	Logger.logger.SetOutput(writer)
	l.output = writer
	l.logger.SetOutput(writer)
}
