package survivalclient

import (
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
)

// GetToken handles access token in the CLI for the whole survival analysis loop
func GetToken(token, username, password string, disableTLSCheck bool) (accessToken string, err error) {

	if len(accessToken) > 0 {
		return
	}
	logrus.Debug("No token provided, requesting token for user ", username, ", disable TLS check: ", disableTLSCheck)
	accessToken, err = utilclient.RetrieveAccessToken(username, password, disableTLSCheck)
	if err != nil {
		return
	}
	return

}
