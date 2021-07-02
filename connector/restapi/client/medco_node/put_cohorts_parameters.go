// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewPutCohortsParams creates a new PutCohortsParams object
// with the default values initialized.
func NewPutCohortsParams() *PutCohortsParams {
	var ()
	return &PutCohortsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPutCohortsParamsWithTimeout creates a new PutCohortsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPutCohortsParamsWithTimeout(timeout time.Duration) *PutCohortsParams {
	var ()
	return &PutCohortsParams{

		timeout: timeout,
	}
}

// NewPutCohortsParamsWithContext creates a new PutCohortsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPutCohortsParamsWithContext(ctx context.Context) *PutCohortsParams {
	var ()
	return &PutCohortsParams{

		Context: ctx,
	}
}

// NewPutCohortsParamsWithHTTPClient creates a new PutCohortsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPutCohortsParamsWithHTTPClient(client *http.Client) *PutCohortsParams {
	var ()
	return &PutCohortsParams{
		HTTPClient: client,
	}
}

/*PutCohortsParams contains all the parameters to send to the API endpoint
for the put cohorts operation typically these are written to a http.Request
*/
type PutCohortsParams struct {

	/*CohortsRequest
	  Cohort that has been updated or created.

	*/
	CohortsRequest PutCohortsBody
	/*Name
	  Name of the cohort to update

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the put cohorts params
func (o *PutCohortsParams) WithTimeout(timeout time.Duration) *PutCohortsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put cohorts params
func (o *PutCohortsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put cohorts params
func (o *PutCohortsParams) WithContext(ctx context.Context) *PutCohortsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put cohorts params
func (o *PutCohortsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put cohorts params
func (o *PutCohortsParams) WithHTTPClient(client *http.Client) *PutCohortsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put cohorts params
func (o *PutCohortsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCohortsRequest adds the cohortsRequest to the put cohorts params
func (o *PutCohortsParams) WithCohortsRequest(cohortsRequest PutCohortsBody) *PutCohortsParams {
	o.SetCohortsRequest(cohortsRequest)
	return o
}

// SetCohortsRequest adds the cohortsRequest to the put cohorts params
func (o *PutCohortsParams) SetCohortsRequest(cohortsRequest PutCohortsBody) {
	o.CohortsRequest = cohortsRequest
}

// WithName adds the name to the put cohorts params
func (o *PutCohortsParams) WithName(name string) *PutCohortsParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the put cohorts params
func (o *PutCohortsParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *PutCohortsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.CohortsRequest); err != nil {
		return err
	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
