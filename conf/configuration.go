package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type flagOptions struct {
	BackupPath string
	ConfigFile string
	DebugMode  bool
	LogLevel   string
	LogFile    string
	Mock       bool
}

// TomlConfig defines TOML file structure
type TomlConfig struct {
	Connection ConnectionInfo
	Backup     BackupInfo
}

// BackupInfo defines the config file backup section
type BackupInfo struct {
	BackupPath string `toml:"backup_path"`
	ConfigFile string `toml:"config_file"`
	DebugMode  string `toml:"debug_mode"`
	LogLevel   string `toml:"log_level"`
	LogFile    string `toml:"log_file"`
	Mock       string `toml:"mock"`
}

// ConnectionInfo defines the config file connection section
type ConnectionInfo struct {
	NifiHost string `toml:"nifi_host"`
	NifiPort string `toml:"nifi_port"`
	NifiUser string `toml:"nifi_user"`
	NifiPass string `toml:"nifi_pass"`
}

// Configurator is an interface for configuration related use.
type Configurator interface {
	Read()
	Write()
}

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
  --mock						        Masquerade as Nifi Server - only used when in debug mode.
  --config <configFile>     TOML format config file to use. [default: hadoop-ottom8r.toml]
  --loglevel <level>        Set loglevel of program. [default: error]
  --logfile <file>          Set file for logging. [default: hadoop-ottom8r.log]`
	return usage
}

// Read info from associated config file
func (tc *TomlConfig) Read() {
	_, err := os.Stat(tc.Backup.ConfigFile)
	if err != nil {
		log.Fatal("Config file is missing: ", tc.Backup.ConfigFile)
	}
	if _, err := toml.DecodeFile(tc.Backup.ConfigFile, tc); err != nil {
		log.Fatal(err)
	}
	log.Debug(fmt.Sprint(tc))
}

// Write out new TOML config file
func (tc *TomlConfig) Write() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(tc); err != nil {
		log.Fatal(err)
	}
	log.Debug(buf.String())
	err := ioutil.WriteFile(tc.Backup.ConfigFile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Wrote new TOML config file.")
}

func getNifiHost(tc TomlConfig) string {
	return tc.Connection.NifiHost
}

func getNifiPort(tc TomlConfig) string {
	return tc.Connection.NifiPort
}

func getNifiUser(tc TomlConfig) string {
	return tc.Connection.NifiUser
}

func getNifiPass(tc TomlConfig) string {
	return tc.Connection.NifiPass
}

func getFlags(arguments map[string]interface{}) *flagOptions {
	flags := &flagOptions{
		ConfigFile: arguments["--config"].(string),
		DebugMode:  arguments["--debug"].(bool),
		LogLevel:   arguments["--loglevel"].(string),
		LogFile:    arguments["--logfile"].(string),
		Mock:       arguments["--mock"].(bool),
	}
	initLogging(flags)
	log.Debug(fmt.Sprint(flags))
	return flags
}
