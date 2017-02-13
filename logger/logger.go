package logger

import (
	"io"
	"os"

	logging "github.com/op/go-logging"
	"github.com/ottom8/hadoop-ottom8r/conf"
)

var log = logging.MustGetLogger("hadoop-ottom8r")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	// Temporary log setting until config is read
	SetupLogger(os.Stdout, logging.ERROR)
}

func GetLogHandle() *logging.Logger {
	return log
}

func InitLogger(flags *conf.FlagOptions) {
	if flags.DebugMode {
		SetupLogger(os.Stdout, logging.DEBUG)
	} else {
		out := setLogFile(flags.LogFile)
		level := readLogLevel(flags.LogLevel)
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
