package nifi

import (
	"github.com/go-resty/resty"
)

// NifiClient internally holds the nifi client state
var nifiClient *resty.Client

func init() {
	nifiClient = resty.New()
}

// SetCredentials sets the user and pass for client
func SetCredentials(user string, password string) {
	nifiClient.SetBasicAuth(user, password)
}
