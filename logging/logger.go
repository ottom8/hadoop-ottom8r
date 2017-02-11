package logging

import (
	"io"
	"os"

	log "github.com/op/go-logging"
)

var logger = log.MustGetLogger("hadoop-ottom8r")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = log.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func initLogger(flags *flagOptions) {
	if flags.DebugMode {
		setupLogger(os.Stdout, log.DEBUG)
	} else {
		out := setLogFile(flags.LogFile)
		level := readLogLevel(flags.LogLevel)
		setupLogger(out, level)
	}
}

func setLogFile(logfile string) *os.File {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("Failed to open log file ", logfile, ":", err)
	}
	return file
}

func readLogLevel(loglevel string) log.Level {
	var level log.Level
	switch loglevel {
	case "debug":
		level = log.DEBUG
	case "info":
		level = log.INFO
	case "notice":
		level = log.NOTICE
	case "warning":
		level = log.WARNING
	case "error":
		level = log.ERROR
	case "critical":
		level = log.INFO
	default:
		log.Fatal("Illegal loglevel provided! Must be one of:" +
			" debug, info, notice, warning, error, critical")
		os.Exit(1)
	}
	return level
}

func setupLogger(out io.Writer, level log.Level) {
	// Initialize log backend
	stdoutBackend := log.NewLogBackend(out, "", 0)
	stdoutBackendFormatter := log.NewBackendFormatter(stdoutBackend, format)
	backendLeveled := log.AddModuleLevel(stdoutBackendFormatter)
	backendLeveled.SetLevel(level, "")
	log.SetBackend(backendLeveled)
}
