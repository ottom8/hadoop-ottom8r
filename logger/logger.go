package logger

import (
	"io"
	"os"

	logging "github.com/op/go-logging"
	"github.com/Sirupsen/logrus"
)

var Log = logging.MustGetLogger("hadoop-ottom8r")
var LogNew = logrus.New()

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//func init() {
//	// Temporary log setting until config is read
//	SetupBareLogger()
//}

func GetLogHandle() *logging.Logger {
	return Log
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
		Log.Fatal("Failed to open log file ", logfile, ":", err)
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
		Log.Fatal("Illegal loglevel provided! Must be one of:" +
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
