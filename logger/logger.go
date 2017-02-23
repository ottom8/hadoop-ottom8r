package logger

import (
	"io"
	"os"

	logging "github.com/op/go-logging"
	"reflect"
	"strconv"
)

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

var log = logging.MustGetLogger("hadoop-ottom8r")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var formatStr = `%{color}%{time:2006-01-02T15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`
var format = logging.MustStringFormatter(formatStr)

//func init() {
//	// Temporary log setting until config is read
//	SetupBareLogger()
//}

func Fatal(out string) {
	log.Fatal(out)
}

func Error(out string) {
	log.Error(out)
}

func Info(out string) {
	log.Info(out)
}

func Debug(out string) {
	log.Debug(out)
}

func OutputStruct(myStruct interface{}) string {
	var out string
	v := reflect.ValueOf(myStruct)

	for i := 0; i < v.NumField(); i++ {
		out += v.Type().Field(i).Name + ":"
		if v.Type().Field(i).Type.String() == "bool" {
			out += strconv.FormatBool(v.Field(i).Bool()) + " "
		} else {
			out += v.Field(i).String() + " "
		}
	}
	return out
}

func InitLogger(debug bool, logFile string, logLevel string) {
	if debug {
		SetupLogger(os.Stdout, logging.DEBUG)
	} else {
		out := setLogFile(logFile)
		level := readLogLevel(logLevel)
		SetupLogger(out, level)
	}
}

func setLogFile(logfile string) *os.File {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file ", logfile, ":", err)
	}
	return file
}

func readLogLevel(loglevel string) logging.Level {
	var level logging.Level
	switch loglevel {
	case "debug":
		level = logging.DEBUG
	case "info":
		level = logging.INFO
	case "notice":
		level = logging.NOTICE
	case "warning":
		level = logging.WARNING
	case "error":
		level = logging.ERROR
	case "critical":
		level = logging.INFO
	default:
		log.Fatal("Illegal loglevel provided! Must be one of:" +
			" debug, info, notice, warning, error, critical")
		os.Exit(1)
	}
	return level
}

func SetupLogger(out io.Writer, level logging.Level) {
	// Initialize log backend
	stdoutBackend := logging.NewLogBackend(out, "", 0)
	stdoutBackendFormatter := logging.NewBackendFormatter(stdoutBackend, format)
	backendLeveled := logging.AddModuleLevel(stdoutBackendFormatter)
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
}

// SetupBareLogger is a convenience function for starting logging before
// the logging configuration is known.
func SetupBareLogger() {
	SetupLogger(os.Stdout, logging.ERROR)
}
