package nifi

import (
	"crypto/tls"
	"time"

	"github.com/ottom8/hadoop-ottom8r/conf"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"fmt"
)

// Currently hold Process Group ID in a const to simplify modification when
// multiple process groups used.
const processGroupId = "8d1d9a24-0158-1000-da01-5d861f4c2834"

// DoBackup contains the business logic for performing Nifi backups.
func DoBackup(tc *conf.TomlConfig) {
	logger.Debug(fmt.Sprint(tc))
	initClient(tc)
	myResp := Call(restHandler(getProcessGroupFlow),
		Request{Id: processGroupId})
	jsonSnippet := ProcessGetProcessGroupFlow(myResp.Body())
	myResp = Call(restHandler(postSnippets),
		Request{Body: jsonSnippet})
	tmplBody := generateBackupName(tc.Connection.NifiHost, processGroupId)
	tmplBody["snippetId"] = ProcessSnippetResponse(myResp.Body())
	myResp = Call(restHandler(postProcessGroupTemplate),
		Request{Id: processGroupId, Body: ProcessTemplateRequest(tmplBody)})
	logger.Debug(fmt.Sprint(myResp))
}

func initClient(tc *conf.TomlConfig) {
	setDefaultHeader(typeAppJson)
	setDefaultHostURL(tc.GetNifiHost())
	//setDefaultRootCert(tc.GetNifiCert())
	setDefaultTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
	token := setCredentials(tc.GetNifiUser(),
		tc.GetNifiPass(), tc.GetNifiToken())
	tc.SetNifiToken(token)
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