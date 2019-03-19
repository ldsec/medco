package i2b2

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
)

// todo doc
// todo: include mechanism to log in and out messages (with configuration)
func i2b2XMLRequest(url string, xmlRequest interface{}, xmlResponse *Response) error {
	marshaledRequest, err := xml.Marshal(xmlRequest)
	if err != nil {
		log.Print("error in i2b2 request marshalling:", err)
		return err
	}

	marshaledRequest = append([]byte(xml.Header), marshaledRequest...)
	httpResponse, err := http.Post(url, "application/xml", bytes.NewBuffer(marshaledRequest))
	if err != nil {
		log.Print("error in i2b2 request HTTP POST:", err)
		return err
	}
	defer httpResponse.Body.Close()

	err = xml.NewDecoder(httpResponse.Body).Decode(xmlResponse)
	if err != nil {
		log.Print("error in i2b2 response unmarshalling:", err)
		return err
	}

	return nil
}
