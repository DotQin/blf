package log

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var (
	log    = logging.MustGetLogger("blf")
	format = logging.MustStringFormatter(
		"BLF %{time:2006-01-02 15:04:05.000} %{color}[%{level:.4s}]%{color:reset} %{message}",
	)
)

func Debug(args ...interface{}) { // lv5
	log.Debug(args...)
}

func Debugf(message string, args ...interface{}) {
	log.Debugf(message, args...)
}

func Info(args ...interface{}) { // lv4
	log.Info(args...)
}

func Infof(message string, args ...interface{}) {
	log.Infof(message, args...)
}

func Notice(args ...interface{}) { // lv3
	log.Notice(args...)
}

func Noticef(message string, args ...interface{}) {
	log.Noticef(message, args...)
}

func Warning(args ...interface{}) { // lv2
	log.Warning(args...)
}

func Warningf(message string, args ...interface{}) {
	log.Warningf(message, args...)
}

func Error(args ...interface{}) { // lv1
	log.Error(args...)
}

func Errorf(message string, args ...interface{}) {
	log.Errorf(message, args...)
}

func GetLevel() string {
	return fmt.Sprintf("%v", logging.GetLevel("blf"))
}

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	log.SetBackend(logging.AddModuleLevel(formatter))
	//logging.SetBackend(formatter) //全局
}
