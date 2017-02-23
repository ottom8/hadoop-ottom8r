package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

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
	DebugMode  bool   `toml:"debug_mode"`
	LogLevel   string `toml:"log_level"`
	LogFile    string `toml:"log_file"`
	Mock       bool   `toml:"mock"`
}

// ConnectionInfo defines the config file connection section
type ConnectionInfo struct {
	NifiHost string 	`toml:"nifi_host"`
	NifiUser string 	`toml:"nifi_user"`
	NifiPass string 	`toml:"nifi_pass"`
	NifiCert string 	`toml:"nifi_cert"`
}

// Configurator is an interface for configuration related use.
type Configurator interface {
	Read()
	Write()
}

//// String method returns the flagOptions object as a string.
//func (fo *flagOptions) String() string {
//	return fmt.Sprintf("%+v", fo)
//}

// String method returns the TomlConfig object as a string.
func (tc *TomlConfig) String() string {
	return "TomlConfig: {" + tc.Backup.String() + " " + tc.Connection.String() + "}"
}

// String method returns the BackupInfo object as a string.
func (bi BackupInfo) String() string {
	return "Backup: {" + logger.OutputStruct(bi) + "}"
}

// String method returns the ConnectionInfo object as a string.
func (ci ConnectionInfo) String() string {
	return fmt.Sprintf("Connection: {NifiHost:%s NifiCert:%s NifiUser:%s NifiPass:********}",
		ci.NifiHost, ci.NifiCert, ci.NifiUser)
}

// Read info from associated config file
func (tc *TomlConfig) Read() {
	_, err := os.Stat(tc.Backup.ConfigFile)
	if err != nil {
		logger.Fatal(fmt.Sprint("Config file is missing: ", tc.Backup.ConfigFile))
	}
	if _, err := toml.DecodeFile(tc.Backup.ConfigFile, tc); err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Debug(fmt.Sprint(tc))
}

// Write out new TOML config file
func (tc *TomlConfig) Write() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(tc); err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Debug(buf.String())
	err := ioutil.WriteFile(tc.Backup.ConfigFile, buf.Bytes(), 0644)
	if err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Info("Wrote new TOML config file.")
}

// GetNifiHost returns the NifiHost config
func (tc *TomlConfig) GetNifiHost() string {
	return tc.Connection.NifiHost
}

// GetNifiUser returns the NifiUser config
func (tc *TomlConfig) GetNifiUser() string {
	return tc.Connection.NifiUser
}

// GetNifiPass returns the NifiPass config
func (tc *TomlConfig) GetNifiPass() string {
	return tc.Connection.NifiPass
}

// GetNifiCert returns the NifiCert config
func (tc *TomlConfig) GetNifiCert() string {
	return tc.Connection.NifiCert
}

func GetFlags(arguments map[string]interface{}) *flagOptions {
	flags := &flagOptions{
		ConfigFile: arguments["--config"].(string),
		DebugMode:  arguments["--debug"].(bool),
		LogLevel:   arguments["--loglevel"].(string),
		LogFile:    arguments["--logfile"].(string),
		Mock:       arguments["--mock"].(bool),
	}
	return flags
}
