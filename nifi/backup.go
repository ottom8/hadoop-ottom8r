package nifi

import (
	"crypto/tls"
	"time"

	"github.com/ottom8/hadoop-ottom8r/conf"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"fmt"
	"io/ioutil"
	"github.com/ottom8/hadoop-ottom8r/util"
	"os"
)

// Currently hold Process Group ID in a const to simplify modification when
// multiple process groups used.
const processGroupId = "8d1d9a24-0158-1000-da01-5d861f4c2834"

// DoBackup contains the business logic for performing Nifi backups.
func DoBackup(ac *conf.AppConfig) {
	var backupMessage string
	logger.Debug(fmt.Sprint(ac))
	initClient(ac)
	myResp := Call(restHandler(getProcessGroupFlow),
		Request{Id: processGroupId})
	jsonSnippet := ProcessGetProcessGroupFlow(myResp.Body())
	myResp = Call(restHandler(postSnippets),
		Request{Body: jsonSnippet})
	tmplBody := generateBackupName(ac.Connection.NifiHost, processGroupId)
	tmplBody["snippetId"] = ProcessSnippetResponse(myResp.Body())
	myResp = Call(restHandler(postProcessGroupTemplate),
		Request{Id: processGroupId, Body: ProcessTemplateRequest(tmplBody)})
	tmplId := ProcessTemplateResponse(myResp.Body())
	myResp = Call(restHandler(getTemplate),
		Request{Id: tmplId})
	logger.Debug(fmt.Sprint(myResp))
	if writeBackup(ac.Backup.BackupPath, tmplBody["name"].(string),
		myResp.Body(), true) {
		backupMessage = fmt.Sprintf("Nifi Backup successful\n")
	} else {
		backupMessage = fmt.Sprintf("Nifi Backup failed\n")
	}

	fmt.Printf(backupMessage)
	logger.Info(backupMessage)
}

func initClient(ac *conf.AppConfig) {
	setDefaultHeader(typeAppJson)
	setDefaultHostURL(ac.GetNifiHost())
	//setDefaultRootCert(ac.GetNifiCert())
	setDefaultTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	//logger.Debug(fmt.Sprintf("%+v", ac.GetNifiPass()))
	token := setCredentials(ac.GetNifiUser(),
		ac.GetNifiPass(), ac.GetNifiToken())
	ac.SetNifiToken(token)
}

func generateBackupName(hostName string, pgId string) map[string]interface{} {
	t := time.Now()
	timeStr := t.Format(time.RFC3339)
	backupName := fmt.Sprintf("backup_nifi_%s_%s", pgId, timeStr)
	backupDescription := fmt.
		Sprintf("Exported template of ProcessGroupId %s from host %s.",
			pgId, hostName)
	return map[string]interface{} {"name": backupName, "description": backupDescription}
}

func writeBackup(backupPath string, backupName string, payload []byte,
		compress bool) bool {
	backupFile := fmt.Sprintf("%s/%s.xml", backupPath, backupName)
	err := ioutil.WriteFile(backupFile, payload, 0644)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if compress {
		util.Tar(backupFile, backupPath)
		util.Gzip(backupFile + ".tar", backupPath)
		os.Remove(backupFile)
		os.Remove(backupFile + ".tar")
	}
	return true
}