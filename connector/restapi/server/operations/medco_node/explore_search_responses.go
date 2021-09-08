// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ExploreSearchOKCode is the HTTP code returned for type ExploreSearchOK
const ExploreSearchOKCode int = 200

/*ExploreSearchOK MedCo-Explore search response.

swagger:response exploreSearchOK
*/
type ExploreSearchOK struct {

	/*
	  In: Body
	*/
	Payload *ExploreSearchOKBody `json:"body,omitempty"`
}

// NewExploreSearchOK creates ExploreSearchOK with default headers values
func NewExploreSearchOK() *ExploreSearchOK {

	return &ExploreSearchOK{}
}

// WithPayload adds the payload to the explore search o k response
func (o *ExploreSearchOK) WithPayload(payload *ExploreSearchOKBody) *ExploreSearchOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore search o k response
func (o *ExploreSearchOK) SetPayload(payload *ExploreSearchOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreSearchOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*ExploreSearchDefault Error response.

swagger:response exploreSearchDefault
*/
type ExploreSearchDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *ExploreSearchDefaultBody `json:"body,omitempty"`
}

// NewExploreSearchDefault creates ExploreSearchDefault with default headers values
func NewExploreSearchDefault(code int) *ExploreSearchDefault {
	if code <= 0 {
		code = 500
	}

	return &ExploreSearchDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the explore search default response
func (o *ExploreSearchDefault) WithStatusCode(code int) *ExploreSearchDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the explore search default response
func (o *ExploreSearchDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the explore search default response
func (o *ExploreSearchDefault) WithPayload(payload *ExploreSearchDefaultBody) *ExploreSearchDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the explore search default response
func (o *ExploreSearchDefault) SetPayload(payload *ExploreSearchDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExploreSearchDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}