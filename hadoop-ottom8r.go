package main

// TODO: Refactor with Cobra and Viper

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	logging "github.com/op/go-logging"
	"github.com/docopt/docopt-go"
	"github.com/ottom8/hadoop-ottom8r/nifi"
	"github.com/ottom8/hadoop-ottom8r/conf"
	"github.com/ottom8/hadoop-ottom8r/logger"
)

func usageOutput() string {
	usage := `hadoop-ottom8r

Usage:
  hadoop-ottom8r [--config <configFile>] [--loglevel <level>] [--logfile <file>] [--debug] [--mock]
  hadoop-ottom8r (-h | --help)
  hadoop-ottom8r --version

Options:
  -h --help                 Show this screen.
  --version                 Show version.
  --debug                   Set loglevel to debug and log to stdout
  --mock					Masquerade as Nifi Server - only used when in debug mode.
  --config <configFile>     TOML format config file to use. [default: hadoop-ottom8r.toml]
  --loglevel <level>        Set loglevel of program. [default: error]
  --logfile <file>          Set file for logging. [default: hadoop-ottom8r.log]`
	return usage
}

func main() {
	// Temporary log setting until config is read
	logger.SetupLogger(os.Stdout, logging.ERROR)
	arguments, _ := docopt.Parse(usageOutput(), nil, true, "hadoop-ottom8r 0.1.0", false)

	utilConfig := conf.LoadConfig(arguments)
	log := logger.GetLogHandle()

	runtime.GOMAXPROCS(4)
	startupMessage := fmt.Sprintf("hadoop-ottom8r started\n")
	fmt.Printf(startupMessage)
	log.Info(startupMessage)

	nifi.DoBackup()
}
