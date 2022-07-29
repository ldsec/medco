// Code generated by go-swagger; DO NOT EDIT.

package explore_statistics

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ExploreStatisticsOKCode is the HTTP code returned for type ExploreStatisticsOK
const ExploreStatisticsOKCode int = 200

/*ExploreStatisticsOK Explore statistics histograms

swagger:response exploreStatisticsOK
*/
type ExploreStatisticsOK struct {

	/*
	  In: Body
	*/
	Payload *ExploreStatisticsOKBody `json:"body,omitempty"`
}

// NewExploreStatisticsOK creates ExploreStatisticsOK with default headers values
func NewExploreStatisticsOK() *ExploreStatisticsOK {

	return &ExploreStatisticsOK{}
}

// WithPayload adds the payload to the explore statistics o k response
func (o *ExploreStatisticsOK) WithPayload(payload *ExploreStatisticsOKBody) *ExploreStatisticsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore statistics o k response
func (o *ExploreStatisticsOK) SetPayload(payload *ExploreStatisticsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreStatisticsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExploreStatisticsBadRequestCode is the HTTP code returned for type ExploreStatisticsBadRequest
const ExploreStatisticsBadRequestCode int = 400

/*ExploreStatisticsBadRequest Bad user input in request.

swagger:response exploreStatisticsBadRequest
*/
type ExploreStatisticsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *ExploreStatisticsBadRequestBody `json:"body,omitempty"`
}

// NewExploreStatisticsBadRequest creates ExploreStatisticsBadRequest with default headers values
func NewExploreStatisticsBadRequest() *ExploreStatisticsBadRequest {

	return &ExploreStatisticsBadRequest{}
}

// WithPayload adds the payload to the explore statistics bad request response
func (o *ExploreStatisticsBadRequest) WithPayload(payload *ExploreStatisticsBadRequestBody) *ExploreStatisticsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore statistics bad request response
func (o *ExploreStatisticsBadRequest) SetPayload(payload *ExploreStatisticsBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreStatisticsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExploreStatisticsNotFoundCode is the HTTP code returned for type ExploreStatisticsNotFound
const ExploreStatisticsNotFoundCode int = 404

/*ExploreStatisticsNotFound Not found.

swagger:response exploreStatisticsNotFound
*/
type ExploreStatisticsNotFound struct {

	/*
	  In: Body
	*/
	Payload *ExploreStatisticsNotFoundBody `json:"body,omitempty"`
}

// NewExploreStatisticsNotFound creates ExploreStatisticsNotFound with default headers values
func NewExploreStatisticsNotFound() *ExploreStatisticsNotFound {

	return &ExploreStatisticsNotFound{}
}

// WithPayload adds the payload to the explore statistics not found response
func (o *ExploreStatisticsNotFound) WithPayload(payload *ExploreStatisticsNotFoundBody) *ExploreStatisticsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore statistics not found response
func (o *ExploreStatisticsNotFound) SetPayload(payload *ExploreStatisticsNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreStatisticsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ExploreStatisticsDefault Error response.

swagger:response exploreStatisticsDefault
*/
type ExploreStatisticsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *ExploreStatisticsDefaultBody `json:"body,omitempty"`
}

// NewExploreStatisticsDefault creates ExploreStatisticsDefault with default headers values
func NewExploreStatisticsDefault(code int) *ExploreStatisticsDefault {
	if code <= 0 {
		code = 500
	}

	return &ExploreStatisticsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the explore statistics default response
func (o *ExploreStatisticsDefault) WithStatusCode(code int) *ExploreStatisticsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the explore statistics default response
func (o *ExploreStatisticsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the explore statistics default response
func (o *ExploreStatisticsDefault) WithPayload(payload *ExploreStatisticsDefaultBody) *ExploreStatisticsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore statistics default response
func (o *ExploreStatisticsDefault) SetPayload(payload *ExploreStatisticsDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreStatisticsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
