// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostCohortsPatientListOKCode is the HTTP code returned for type PostCohortsPatientListOK
const PostCohortsPatientListOKCode int = 200

/*PostCohortsPatientListOK Queried patient list

swagger:response postCohortsPatientListOK
*/
type PostCohortsPatientListOK struct {

	/*
	  In: Body
	*/
	Payload *PostCohortsPatientListOKBody `json:"body,omitempty"`
}

// NewPostCohortsPatientListOK creates PostCohortsPatientListOK with default headers values
func NewPostCohortsPatientListOK() *PostCohortsPatientListOK {

	return &PostCohortsPatientListOK{}
}

// WithPayload adds the payload to the post cohorts patient list o k response
func (o *PostCohortsPatientListOK) WithPayload(payload *PostCohortsPatientListOKBody) *PostCohortsPatientListOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post cohorts patient list o k response
func (o *PostCohortsPatientListOK) SetPayload(payload *PostCohortsPatientListOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCohortsPatientListOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostCohortsPatientListForbiddenCode is the HTTP code returned for type PostCohortsPatientListForbidden
const PostCohortsPatientListForbiddenCode int = 403

/*PostCohortsPatientListForbidden Request is valid and user is authenticated, but not authorized to perform this action.

swagger:response postCohortsPatientListForbidden
*/
type PostCohortsPatientListForbidden struct {

	/*
	  In: Body
	*/
	Payload *PostCohortsPatientListForbiddenBody `json:"body,omitempty"`
}

// NewPostCohortsPatientListForbidden creates PostCohortsPatientListForbidden with default headers values
func NewPostCohortsPatientListForbidden() *PostCohortsPatientListForbidden {

	return &PostCohortsPatientListForbidden{}
}

// WithPayload adds the payload to the post cohorts patient list forbidden response
func (o *PostCohortsPatientListForbidden) WithPayload(payload *PostCohortsPatientListForbiddenBody) *PostCohortsPatientListForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post cohorts patient list forbidden response
func (o *PostCohortsPatientListForbidden) SetPayload(payload *PostCohortsPatientListForbiddenBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCohortsPatientListForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostCohortsPatientListNotFoundCode is the HTTP code returned for type PostCohortsPatientListNotFound
const PostCohortsPatientListNotFoundCode int = 404

/*PostCohortsPatientListNotFound Not found.

swagger:response postCohortsPatientListNotFound
*/
type PostCohortsPatientListNotFound struct {

	/*
	  In: Body
	*/
	Payload *PostCohortsPatientListNotFoundBody `json:"body,omitempty"`
}

// NewPostCohortsPatientListNotFound creates PostCohortsPatientListNotFound with default headers values
func NewPostCohortsPatientListNotFound() *PostCohortsPatientListNotFound {

	return &PostCohortsPatientListNotFound{}
}

// WithPayload adds the payload to the post cohorts patient list not found response
func (o *PostCohortsPatientListNotFound) WithPayload(payload *PostCohortsPatientListNotFoundBody) *PostCohortsPatientListNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post cohorts patient list not found response
func (o *PostCohortsPatientListNotFound) SetPayload(payload *PostCohortsPatientListNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCohortsPatientListNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*PostCohortsPatientListDefault Error response.

swagger:response postCohortsPatientListDefault
*/
type PostCohortsPatientListDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *PostCohortsPatientListDefaultBody `json:"body,omitempty"`
}

// NewPostCohortsPatientListDefault creates PostCohortsPatientListDefault with default headers values
func NewPostCohortsPatientListDefault(code int) *PostCohortsPatientListDefault {
	if code <= 0 {
		code = 500
	}

	return &PostCohortsPatientListDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post cohorts patient list default response
func (o *PostCohortsPatientListDefault) WithStatusCode(code int) *PostCohortsPatientListDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post cohorts patient list default response
func (o *PostCohortsPatientListDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post cohorts patient list default response
func (o *PostCohortsPatientListDefault) WithPayload(payload *PostCohortsPatientListDefaultBody) *PostCohortsPatientListDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post cohorts patient list default response
func (o *PostCohortsPatientListDefault) SetPayload(payload *PostCohortsPatientListDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostCohortsPatientListDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
