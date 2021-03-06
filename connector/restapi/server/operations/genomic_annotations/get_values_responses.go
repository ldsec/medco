// Code generated by go-swagger; DO NOT EDIT.

package genomic_annotations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetValuesOKCode is the HTTP code returned for type GetValuesOK
const GetValuesOKCode int = 200

/*GetValuesOK Queried annotation values.

swagger:response getValuesOK
*/
type GetValuesOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewGetValuesOK creates GetValuesOK with default headers values
func NewGetValuesOK() *GetValuesOK {

	return &GetValuesOK{}
}

// WithPayload adds the payload to the get values o k response
func (o *GetValuesOK) WithPayload(payload []string) *GetValuesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get values o k response
func (o *GetValuesOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetValuesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]string, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetValuesNotFoundCode is the HTTP code returned for type GetValuesNotFound
const GetValuesNotFoundCode int = 404

/*GetValuesNotFound Annotation not found.

swagger:response getValuesNotFound
*/
type GetValuesNotFound struct {
}

// NewGetValuesNotFound creates GetValuesNotFound with default headers values
func NewGetValuesNotFound() *GetValuesNotFound {

	return &GetValuesNotFound{}
}

// WriteResponse to the client
func (o *GetValuesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

/*GetValuesDefault Error response.

swagger:response getValuesDefault
*/
type GetValuesDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *GetValuesDefaultBody `json:"body,omitempty"`
}

// NewGetValuesDefault creates GetValuesDefault with default headers values
func NewGetValuesDefault(code int) *GetValuesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetValuesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get values default response
func (o *GetValuesDefault) WithStatusCode(code int) *GetValuesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get values default response
func (o *GetValuesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get values default response
func (o *GetValuesDefault) WithPayload(payload *GetValuesDefaultBody) *GetValuesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get values default response
func (o *GetValuesDefault) SetPayload(payload *GetValuesDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetValuesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
