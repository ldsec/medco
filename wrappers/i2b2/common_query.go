package i2b2

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// todo doc
func i2b2XMLRequest(url string, xmlRequest interface{}, xmlResponse *Response) error {
	logrus.Info("i2b2 XML request to " + url)

	// marshal request
	marshaledRequest, err := xml.MarshalIndent(xmlRequest, "  ", "    ")
	if err != nil {
		logrus.Error("error in i2b2 request marshalling:", err)
		return err
	}
	marshaledRequest = append([]byte(xml.Header), marshaledRequest...)
	logrus.Debug("i2b2 request:\n", string(marshaledRequest))

	// execute HTTP request
	httpResponse, err := http.Post(url, "application/xml", bytes.NewBuffer(marshaledRequest))
	if err != nil {
		logrus.Error("error in i2b2 request HTTP POST:", err)
		return err
	}
	defer httpResponse.Body.Close()

	// unmarshal response
	httpBody, err := ioutil.ReadAll(httpResponse.Body)
	logrus.Debug("i2b2 response:\n", string(httpBody))

	err = xml.Unmarshal(httpBody, xmlResponse)
	if err != nil {
		logrus.Error("error in i2b2 response unmarshalling:", err)
		return err
	}

	// check i2b2 request status
	err = xmlResponse.checkStatus()
	if err != nil {
		logrus.Error("i2b2 error:", err)
		return err
	}

	return nil
}
