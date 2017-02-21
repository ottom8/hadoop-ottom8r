package nifi

import (
	"fmt"

	"github.com/ottom8/hadoop-ottom8r/conf"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"crypto/tls"
)

// DoBackup contains the business logic for performing Nifi backups
func DoBackup(tc *conf.TomlConfig) {
	//output := fmt.Sprintf("%+v", tc)
	//logger.Log.Debug(output)
	logger.Log.Debug(tc.Inspect())
	initClient(tc)
	myResp := Call(restHandler(getProcessGroup), Request{Id: "root", Body: ""})
	logger.Log.Debug(fmt.Sprintf("%+v", myResp))
}

func initClient(tc *conf.TomlConfig) {
	setCredentials(tc.GetNifiUser(), tc.GetNifiPass())
	setDefaultHeader("application/json")
	setDefaultHostURL(tc.GetNifiHost())
	//setDefaultRootCert(tc.GetNifiCert())
	setDefaultTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
}