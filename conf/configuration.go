package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/ottom8/hadoop-ottom8r/logger"
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

func GetNifiUser(tc TomlConfig) string {
	return tc.Connection.NifiUser
}

// GetNifiPass returns the NifiPass config
func GetNifiPass(tc TomlConfig) string {
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
	return flags
}

func LoadConfig(arguments map[string]interface{}) *TomlConfig {
	// Initialize flags instance variable
	flags := getFlags(arguments)
	logger.InitLogger(flags)
	log := logger.GetLogHandle()
	log.Debug(fmt.Sprint(flags))
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
	return utilConfig
}
