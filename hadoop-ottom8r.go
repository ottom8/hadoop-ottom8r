package main

// TODO: Refactor with Cobra and Viper

import (
	"fmt"
	//"runtime"
	//"strconv"

	logging "github.com/op/go-logging"
	"github.com/docopt/docopt-go"
	"github.com/ottom8/hadoop-ottom8r/conf"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"github.com/ottom8/hadoop-ottom8r/nifi"
	"os"
)

var log *logging.Logger

func usageOutput() string {
	usage := `hadoop-ottom8r

Usage:
  hadoop-ottom8r [--config <configFile>] [--loglevel <level>] [--logfile <file>] [--debug] [--mock] [--encrypt]
  hadoop-ottom8r (-h | --help)
  hadoop-ottom8r --version

Options:
  -h --help                 Show this screen.
  --version                 Show version.
  --debug                   Set loglevel to debug and log to stdout
  --mock                    Masquerade as Nifi Server - only used when in debug mode.
  --encrypt                 Encrypt plaintext passwords in config file and immediately exit.
  --config <configFile>     TOML format config file to use. [default: hadoop-ottom8r.toml]
  --loglevel <level>        Set loglevel of program. [default: error]
  --logfile <file>          Set file for logging. [default: hadoop-ottom8r.log]`
	return usage
}

func main() {
	logger.SetupBareLogger()
	//log = logger.GetLogHandle()
	arguments, optErr := docopt.Parse(usageOutput(), nil, true, "hadoop-ottom8r 0.1.0", false)
	if optErr != nil {
		logger.Fatal(fmt.Sprint(optErr))
	}

	utilConfig := loadConfig(arguments)

	//runtime.GOMAXPROCS(4)
	startupMessage := fmt.Sprintf("hadoop-ottom8r started\n")
	fmt.Printf(startupMessage)
	logger.Info(startupMessage)

	nifi.DoBackup(utilConfig)
}

func loadConfig(arguments map[string]interface{}) *conf.AppConfig {
	// Initialize flags instance variable
	flags := conf.GetFlags(arguments)
	logger.Debug(fmt.Sprint(flags))
	utilConfig := new(conf.AppConfig)
	utilConfig.Backup.ConfigFile = flags.ConfigFile
	//utilConfig.Backup.DebugMode = strconv.FormatBool(flags.DebugMode)
	utilConfig.Backup.DebugMode = flags.DebugMode
	utilConfig.Backup.LogLevel = flags.LogLevel
	utilConfig.Backup.LogFile = flags.LogFile
	//utilConfig.Backup.Mock = strconv.FormatBool(flags.Mock)
	utilConfig.Backup.Mock = flags.Mock
	utilConfig.TomlConfig.Read()
	logger.InitLogger(utilConfig.Backup.DebugMode, utilConfig.Backup.LogFile,
		utilConfig.Backup.LogLevel)
	utilConfig.AuthConfig.Read()
	logger.Debug(fmt.Sprint(utilConfig))
	logger.Info("Loaded configuration.")
	logger.Debug(fmt.Sprint(utilConfig))
	if flags.Encrypt {
		doEncryption(utilConfig)
	}
	return utilConfig
}

func doEncryption(ac *conf.AppConfig) {
	ac.Encrypt()
	os.Exit(0)
}
