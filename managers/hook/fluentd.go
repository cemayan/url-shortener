package hook

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
	"log"
)

// FluentdlogHook to send logs
type FluentdlogHook struct {
}

// SendLogToFluentd send event to fluentd
// entry.Data includes to given fields from developer
func (hook *FluentdlogHook) SendLogToFluentd(entry *logrus.Entry, tag string) {

	logger, err := fluent.New(fluent.Config{})
	if err != nil {
		fmt.Printf("Fluentd connection error : %v\n", err.Error())
	} else {

		entry.Data["message"] = entry.Message

		defer logger.Close()
		error := logger.Post(tag, entry.Data)
		if error != nil {
			log.Fatalln(error.Error())
		}
	}
}

func (hook *FluentdlogHook) Fire(entry *logrus.Entry) error {

	switch entry.Level {
	case logrus.PanicLevel:
		hook.SendLogToFluentd(entry, "panic")
	case logrus.FatalLevel:
		hook.SendLogToFluentd(entry, "fatal")
	case logrus.ErrorLevel:
		hook.SendLogToFluentd(entry, "error")
	default:
		return nil
	}

	return nil
}

func (hook *FluentdlogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewFluentdHook() *FluentdlogHook {
	return &FluentdlogHook{}
}
