package main

// TODO: Refactor with Cobra and Viper

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/ottom8/hadoop-ottom8r/nifi"
)

func main() {
	setupLogger(os.Stdout, logging.ERROR)
	arguments, _ := docopt.Parse(usageOutput(), nil, true, "hadoop-ottom8r 0.1.0", false)

	// Initialize flags instance variable
	flags := getFlags(arguments)
	utilConfig := new(TomlConfig)
	utilConfig.Backup.ConfigFile = flags.ConfigFile
	utilConfig.Backup.DebugMode = strconv.FormatBool(flags.DebugMode)
	utilConfig.Backup.LogLevel = flags.LogLevel
	utilConfig.Backup.LogFile = flags.LogFile
	utilConfig.Backup.Mock = strconv.FormatBool(flags.Mock)
	utilConfig.Read()
	log.Debug(fmt.Sprint(utilConfig))
	log.Info("Loaded configuration.")
	log.Debug(fmt.Sprint(utilConfig))

	runtime.GOMAXPROCS(4)
	startupMessage := fmt.Sprintf("hadoop-ottom8r started\n")
	fmt.Printf(startupMessage)
	log.Info(startupMessage)

	nifi.DoBackup()
}
