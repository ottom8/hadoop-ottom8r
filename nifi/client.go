package nifi

import (
	"github.com/go-resty/resty"
	"fmt"
	"github.com/ottom8/hadoop-ottom8r/logger"
	"crypto/tls"
)

const (
	hdrContentType = "Content-Type"
	hdrAccept = "Accept"
	typeAppJson = "application/json"
	typeText = "text/plain"
	typeForm = "application/x-www-form-urlencoded"

	endpointAccess = "access"
	endpointAccessToken = "access/token"
	endpointFlow = "flow"
	endpointProcessGroups = "process-groups"
)

type CredentialsInfo struct {
	user string
	password string
	authToken string
}

var credentials CredentialsInfo

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
		logger.Error(fmt.Sprintf(err.Error()))
	}
	return resp
}

// NifiClient internally holds the nifi client state
var nifiClient *resty.Client

func init() {
	nifiClient = resty.New()
}

// SetClientCredentials sets the client token (if present) else user and pass
// for client
func setClientCredentials() {
	nifiClient.SetAuthToken(credentials.authToken)
	if !isAuthValid() {
		myResp := Call(restHandler(postAccessToken), Request{})
		setAuthToken(string(myResp.Body()))
	}
}

// SetCredentials sets the package credentials var and establishes client
// auth returns current token
func setCredentials(user string, password string, token string) string {
	credentials = CredentialsInfo{user: user, password: password,
		authToken:token}
	setClientCredentials()
	return credentials.authToken
}

func isAuthValid() bool {
	if credentials.authToken == "" {
		return false
	}
	myResp := Call(restHandler(getAccess), Request{})
	return myResp.StatusCode() == 200
}

func setAuthToken(token string) {
	credentials.authToken = token
	nifiClient.SetAuthToken(token)

}

func setDefaultHeader(contentType string) {
	nifiClient.
		SetHeader(hdrContentType, contentType).
		SetHeader(hdrAccept, contentType)
}

func setDefaultHostURL(hostURL string) {
	nifiClient.SetHostURL(fmt.Sprintf("%s/nifi-api", hostURL))
}

func setDefaultRootCert(pemFilePath string) {
	resp := nifiClient.SetRootCertificate(fmt.Sprintf("resources/%s", pemFilePath))
	logger.Debug(fmt.Sprint(resp))
}

func setDefaultTLSClientConfig(config *tls.Config) {
	resp := nifiClient.SetTLSClientConfig(config)
	logger.Debug(fmt.Sprint(resp))
}

// postAccessToken performs a login against Nifi to create and consume an
// auth token
func postAccessToken(_ Request) (*resty.Response, error) {
	resp, err := nifiClient.R().
		SetHeader(hdrContentType, typeForm).
		SetHeader(hdrAccept, typeText).
		SetBody(fmt.Sprintf("username=%s&password=%s",
			credentials.user, credentials.password)).
		Post(fmt.Sprintf("/%s", endpointAccessToken))
	return resp, err
}

func getAccess(_ Request) (*resty.Response, error) {
	resp, err := nifiClient.R().
			Get(fmt.Sprintf("/%s", endpointAccess))
	return resp, err
}

func getProcessGroup(req Request) (*resty.Response, error) {
	resp, err := nifiClient.R().
		Get(fmt.Sprintf("/%s/%s", endpointProcessGroups, req.Id))
	return resp, err
}

func getProcessGroupFlow(req Request) (*resty.Response, error) {
	resp, err := nifiClient.R().
			Get(fmt.Sprintf("/%s/%s/%s", endpointFlow,
		endpointProcessGroups, req.Id))
	return resp, err
}

func postProcessGroupTemplate(req Request) (*resty.Response, error) {
	logger.Debug(req.Body)
	resp, err := nifiClient.R().
		Post(fmt.Sprintf("/%s/%s/templates", endpointProcessGroups, req.Id))
	return resp, err
}