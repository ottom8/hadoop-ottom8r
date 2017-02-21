package nifi

import (
	"github.com/go-resty/resty"
	"fmt"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"crypto/tls"
)

const (
	endpointProcessGroups = "process-groups"
)

type Request struct {
	Id		string
	Body	string
}

type Handler interface {
	RestCall(Request) *resty.Response
}

type restHandler func(Request) (*resty.Response, error)

func Call(fn Handler, req Request) *resty.Response {
	return Handler(fn).RestCall(req)
}

// RestCall is a wrapper to make REST calls
func (fn restHandler) RestCall(req Request) *resty.Response {
	resp, err := fn(req); if err != nil {
		logger.Log.Error(err.Error())
	}
	return resp
}

// NifiClient internally holds the nifi client state
var nifiClient *resty.Client

func init() {
	nifiClient = resty.New()
}

// SetCredentials sets the user and pass for client
func setCredentials(user string, password string) {
	nifiClient.SetBasicAuth(user, password)
}

func setDefaultHeader(header string) {
	nifiClient.
		SetHeader("Content-Type", header).
		SetHeader("Accept", header)
}

func setDefaultHostURL(hostURL string) {
	nifiClient.SetHostURL(fmt.Sprintf("%s/nifi-api", hostURL))
}

func setDefaultRootCert(pemFilePath string) {
	resp := nifiClient.SetRootCertificate(fmt.Sprintf("resources/%s", pemFilePath))
	logger.Log.Debug(fmt.Sprintf("%+v", resp))
}

func setDefaultTLSClientConfig(config *tls.Config) {
	resp := nifiClient.SetTLSClientConfig(config)
	logger.Log.Debug(fmt.Sprintf("%+v", resp))
}

func getProcessGroup(req Request) (*resty.Response, error) {
	resp, err := nifiClient.R().
		Get(fmt.Sprintf("/%s/%s", endpointProcessGroups, req.Id))
	//if err != nil {
	//	logger.Log.Error(err.Error())
	//}
	return resp, err
}

func postProcessGroupTemplate(req Request) (*resty.Response, error) {
	logger.Log.Debug(fmt.Sprintf("%+v", req.Body))
	resp, err := nifiClient.R().
		Post(fmt.Sprintf("/%s/%s/templates", endpointProcessGroups, req.Id))
	//if err != nil {
	//	logger.Log.Error(err.Error())
	//}
	return resp, err
}