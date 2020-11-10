// Code generated by go-swagger; DO NOT EDIT.

package survival_analysis

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// SurvivalAnalysisOKCode is the HTTP code returned for type SurvivalAnalysisOK
const SurvivalAnalysisOKCode int = 200

/*SurvivalAnalysisOK Queried survival analysis

swagger:response survivalAnalysisOK
*/
type SurvivalAnalysisOK struct {

	/*
	  In: Body
	*/
	Payload *SurvivalAnalysisOKBody `json:"body,omitempty"`
}

// NewSurvivalAnalysisOK creates SurvivalAnalysisOK with default headers values
func NewSurvivalAnalysisOK() *SurvivalAnalysisOK {

	return &SurvivalAnalysisOK{}
}

// WithPayload adds the payload to the survival analysis o k response
func (o *SurvivalAnalysisOK) WithPayload(payload *SurvivalAnalysisOKBody) *SurvivalAnalysisOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the survival analysis o k response
func (o *SurvivalAnalysisOK) SetPayload(payload *SurvivalAnalysisOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SurvivalAnalysisOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SurvivalAnalysisBadRequestCode is the HTTP code returned for type SurvivalAnalysisBadRequest
const SurvivalAnalysisBadRequestCode int = 400

/*SurvivalAnalysisBadRequest Bad user input in request.

swagger:response survivalAnalysisBadRequest
*/
type SurvivalAnalysisBadRequest struct {

	/*
	  In: Body
	*/
	Payload *SurvivalAnalysisBadRequestBody `json:"body,omitempty"`
}

// NewSurvivalAnalysisBadRequest creates SurvivalAnalysisBadRequest with default headers values
func NewSurvivalAnalysisBadRequest() *SurvivalAnalysisBadRequest {

	return &SurvivalAnalysisBadRequest{}
}

// WithPayload adds the payload to the survival analysis bad request response
func (o *SurvivalAnalysisBadRequest) WithPayload(payload *SurvivalAnalysisBadRequestBody) *SurvivalAnalysisBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the survival analysis bad request response
func (o *SurvivalAnalysisBadRequest) SetPayload(payload *SurvivalAnalysisBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SurvivalAnalysisBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SurvivalAnalysisNotFoundCode is the HTTP code returned for type SurvivalAnalysisNotFound
const SurvivalAnalysisNotFoundCode int = 404

/*SurvivalAnalysisNotFound Not found.

swagger:response survivalAnalysisNotFound
*/
type SurvivalAnalysisNotFound struct {

	/*
	  In: Body
	*/
	Payload *SurvivalAnalysisNotFoundBody `json:"body,omitempty"`
}

// NewSurvivalAnalysisNotFound creates SurvivalAnalysisNotFound with default headers values
func NewSurvivalAnalysisNotFound() *SurvivalAnalysisNotFound {

	return &SurvivalAnalysisNotFound{}
}

// WithPayload adds the payload to the survival analysis not found response
func (o *SurvivalAnalysisNotFound) WithPayload(payload *SurvivalAnalysisNotFoundBody) *SurvivalAnalysisNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the survival analysis not found response
func (o *SurvivalAnalysisNotFound) SetPayload(payload *SurvivalAnalysisNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SurvivalAnalysisNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*SurvivalAnalysisDefault Error response.

swagger:response survivalAnalysisDefault
*/
type SurvivalAnalysisDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *SurvivalAnalysisDefaultBody `json:"body,omitempty"`
}

// NewSurvivalAnalysisDefault creates SurvivalAnalysisDefault with default headers values
func NewSurvivalAnalysisDefault(code int) *SurvivalAnalysisDefault {
	if code <= 0 {
		code = 500
	}

	return &SurvivalAnalysisDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the survival analysis default response
func (o *SurvivalAnalysisDefault) WithStatusCode(code int) *SurvivalAnalysisDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the survival analysis default response
func (o *SurvivalAnalysisDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the survival analysis default response
func (o *SurvivalAnalysisDefault) WithPayload(payload *SurvivalAnalysisDefaultBody) *SurvivalAnalysisDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the survival analysis default response
func (o *SurvivalAnalysisDefault) SetPayload(payload *SurvivalAnalysisDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SurvivalAnalysisDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
