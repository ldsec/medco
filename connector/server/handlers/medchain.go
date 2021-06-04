package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/medchain"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

// MedchainWsProxyHandler handles /medchain/ws API endpoint
func MedchainWsProxyHandler(params medchain.WsProxyParams, _ *models.User) middleware.Responder {

	if !utilserver.MedChainEnabled {
		err := fmt.Errorf("request to proxy MedChain request while MedChain support is disabled on this instance")
		logrus.Error(err)
		return medchain.NewWsProxyForbidden().WithPayload(&medchain.WsProxyForbiddenBody{Message: err.Error()})
	}

	medchainNodeURL, err := url.Parse(utilserver.MedChainWsURL)
	if err != nil {
		logrus.Error("parsing websocket URL:", err)
		return medchain.NewWsProxyDefault(500).WithPayload(&medchain.WsProxyDefaultBody{Message: err.Error()})
	}

	// generate the new URL
	proxyBackend := func(r *http.Request) *url.URL {

		// extract the URL path by removing elements before the medco proxy path "/medchain/ws"
		splitURL := strings.Split(r.URL.Path, "medchain/ws")
		if len(splitURL) == 1 {
			medchainNodeURL.Path = "/"
		} else {
			medchainNodeURL.Path = splitURL[1]
		}

		// extract URL Fragment and RawQuery
		medchainNodeURL.Fragment = r.URL.Fragment
		medchainNodeURL.RawQuery = r.URL.RawQuery
		return medchainNodeURL
	}

	wsProxy := websocketproxy.WebsocketProxy{
		Backend: proxyBackend,

		Director: func(incoming *http.Request, out http.Header) {
			// do not keep original host (which is the default proxy behavior)
			out.Set("Host", medchainNodeURL.Host)
		},

		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				// disable the origin check: it is to be done by the reverse proxy before that
				return true
			},
		},
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, _ runtime.Producer) {
		logrus.Infof("proxying websocket connection for MedChain to %v", medchainNodeURL.String())
		wsProxy.ServeHTTP(rw, params.HTTPRequest)
	})
}
