package common

import (
	"context"
	"fmt"
	"net/http"
	"time"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

//for client !!

func (clientSurvivalAnalysis *SurvivalAnalysis) submitToNode() (results []string, err error) {
	//magicNumber
	params := GetSurvivalAnalysisParameter{Command: clientSurvivalAnalysis, timeout: time.Duration(300) * time.Second}
	authInfo := httptransport.BearerToken(clientSurvivalAnalysis.authToken)
	response, err := clientSurvivalAnalysis.httpMedCoClient.Transport.Submit(&runtime.ClientOperation{
		ID:                 "getSurvivalAnalysis",
		Method:             "GET",
		PathPattern:        "/survival-analysis/{granularity}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             &params,
		Reader:             &GetSurvivalAnalysisReader{formats: clientSurvivalAnalysis.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return
	}
	success, ok := response.(*GetSurvivalAnalysisOK)

	if ok {
		results = success.Payload
		return
	}
	notFound, notFoundOk := response.(*GetSurvivalAnalysisNotFound)
	if notFoundOk {
		results = nil
		err = notFound
		return
	}

	// unexpected success response
	unexpectedSuccess := response.(*GetSurvivalAnalysisDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

//to for node submitting
type GetSurvivalAnalysisParameter struct {
	Command    *SurvivalAnalysis
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

func (this *GetSurvivalAnalysisParameter) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	err := r.SetTimeout(this.timeout)
	if err != nil {
		return err
	}

	//var res []error
	//err = r.SetPathParam("survival", this.Command.granularity)

	err = r.SetPathParam("granularity", this.Command.granularity)
	return err

}

type GetSurvivalAnalysisReader struct {
	formats strfmt.Registry
}

func (this *GetSurvivalAnalysisReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := &GetSurvivalAnalysisOK{}
		if err := result.readResponse(response, consumer, this.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := &GetSurvivalAnalysisNotFound{}
		if err := result.readResponse(response, consumer, this.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := &GetSurvivalAnalysisDefault{_statusCode: response.Code()}
		if err := result.readResponse(response, consumer, this.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result

	}
}

type GetSurvivalAnalysisOK struct {
	Payload []string
}

func (this *GetSurvivalAnalysisOK) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{granularity}/][%d] getVariantsOK  %+v", 200, this.Payload)

}

func (this *GetSurvivalAnalysisOK) GetPayload() []string {
	return this.Payload
}

func (this *GetSurvivalAnalysisOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) (err error) {
	err = consumer.Consume(response.Body(), &this.Payload)
	return

}

type GetSurvivalAnalysisNotFound struct{}

func (this *GetSurvivalAnalysisNotFound) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{garnularity}/][%d] getSurvivalAnalysisNotFound", 404)
}

func (this *GetSurvivalAnalysisNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	return nil
}

type DefaultBody struct {
	Message string `json:"message,omitempty"`
}

func (this *DefaultBody) MarshalBinary() (res []byte, err error) {
	if this == nil {
		return
	}
	res, err = swag.WriteJSON(this)
	return

}
func (this *DefaultBody) UnmarshalBinary(b []byte) (err error) {
	var res DefaultBody
	err = swag.ReadJSON(b, &res)
	if err != nil {
		return
	}
	*this = res
	return
}
func (this DefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

type GetSurvivalAnalysisDefault struct {
	_statusCode int
	Payload     *DefaultBody
}

func (this *GetSurvivalAnalysisDefault) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{granularity}][%d] getValues default  %+v", this._statusCode, this.Payload)
}

func (this *GetSurvivalAnalysisDefault) Code() int {
	return this._statusCode
}

func (this *GetSurvivalAnalysisDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) (err error) {
	this.Payload = &DefaultBody{}
	err = consumer.Consume(response.Body(), this.Payload)
	return
}
