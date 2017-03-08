package conf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"github.com/ottom8/hadoop-ottom8r/util"
)

const (
	authConfigFile = "auth.toml"
	encryptKey = "H@d00pS3cret3ncryptK3yS3cur3RlLg"
)

type flagOptions struct {
	BackupPath string
	ConfigFile string
	DebugMode  bool
	LogLevel   string
	LogFile    string
	Mock       bool
	Encrypt    bool
}

// AppConfig defines all configuration for this app
type AppConfig struct {
	AuthConfig
	TomlConfig
}

// AuthConfig defines Auth file structure
type AuthConfig struct {
	NifiToken string `toml:"nifi_token"`
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

// String method returns the AppConfig object as a string.
func (ac *AppConfig) String() string {
	return "AppConfig: {" + ac.TomlConfig.String() + " " + ac.AuthConfig.String() + "}"
}

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

// String method returns the AuthConfig object as a string.
func (ac AuthConfig) String() string {
	return fmt.Sprintf("AuthConfig: {NifiToken:%s}", ac.NifiToken)
}

// Read info from associated config file
func (ac *AuthConfig) Read() {
	_, err := os.Stat(authConfigFile)
	if err != nil {
		logger.Info(fmt.Sprintf("AuthConfig file %s is missing, creating. ",
			authConfigFile))
		ac.Write()
	}
	if _, err := toml.DecodeFile(authConfigFile, ac); err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Debug(fmt.Sprint(ac))
}

// Write out new auth config file
func (ac *AuthConfig) Write() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(ac); err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Debug(buf.String())
	err := ioutil.WriteFile(authConfigFile, buf.Bytes(), 0644)
	if err != nil {
		logger.Fatal(fmt.Sprint(err))
	}
	logger.Info("Wrote new auth config file.")
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

// Encrypt performs encryption on plaintext passwords in config file.
func (tc *TomlConfig) Encrypt() {
	tc.SetNifiPass(tc.Connection.NifiPass)
}

// GetNifiHost returns the NifiHost config
func (tc *TomlConfig) GetNifiHost() string {
	return tc.Connection.NifiHost
}

// GetNifiUser returns the NifiUser config
func (tc *TomlConfig) GetNifiUser() string {
	return tc.Connection.NifiUser
}

// GetNifiPass decrypts and returns the NifiPass config
func (tc *TomlConfig) GetNifiPass() string {
	return util.Base64Decrypt(tc.Connection.NifiPass, encryptKey)
}

// SetNifiPass encrypts and stores the NifiPass config
func (tc *TomlConfig) SetNifiPass(pass string) {
	outbound := util.Base64Encrypt(pass, encryptKey)
	if outbound != tc.Connection.NifiPass {
		tc.Connection.NifiPass = outbound
		tc.Write()
	}
}

// GetNifiCert returns the NifiCert config
func (tc *TomlConfig) GetNifiCert() string {
	return tc.Connection.NifiCert
}

// GetNifiToken returns the NifiToken config
func (ac *AuthConfig) GetNifiToken() string {
	return ac.NifiToken
}

// SetNifiToken updates the NifiToken config
func (ac *AuthConfig) SetNifiToken(token string) {
	if token != ac.NifiToken {
		ac.NifiToken = token
		ac.Write()
	}
}

func GetFlags(arguments map[string]interface{}) *flagOptions {
	flags := &flagOptions{
		ConfigFile: arguments["--config"].(string),
		DebugMode:  arguments["--debug"].(bool),
		LogLevel:   arguments["--loglevel"].(string),
		LogFile:    arguments["--logfile"].(string),
		Mock:       arguments["--mock"].(bool),
		Encrypt:    arguments["--encrypt"].(bool),
	}
	return flags
}
