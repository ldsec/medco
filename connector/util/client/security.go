package utilclient

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

// oidcTokenResp contains the response to an OIDC token request
type oidcTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// RetrieveAccessToken requests JWT from OIDC provider
func RetrieveAccessToken(username string, password string, disableTLSCheck bool) (token string, err error) {

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: disableTLSCheck,
			},
		},
	}

	httpResp, err := httpClient.PostForm(OidcReqTokenURL, url.Values{
		"grant_type": {"password"},
		"client_id":  {OidcReqTokenClientID},
		"username":   {username},
		"password":   {password},
	})

	if err != nil {
		logrus.Error("OIDC request token error: ", err)
		return
	}

	if httpResp.StatusCode != 200 {
		err = errors.New("OIDC request token error (code " + fmt.Sprint(httpResp.StatusCode) + ")")
		logrus.Error(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	parsedResp := &oidcTokenResp{}
	err = json.Unmarshal(bodyBytes, parsedResp)
	if err != nil {
		logrus.Error("OIDC request token error unmarshalling: ", err)
		return
	}

	logrus.Info("OIDC request token successfully authenticated")
	logrus.Debug("OIDC request token: " + parsedResp.AccessToken)
	return parsedResp.AccessToken, nil
}
