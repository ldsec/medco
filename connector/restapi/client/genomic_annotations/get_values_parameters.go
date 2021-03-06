// Code generated by go-swagger; DO NOT EDIT.

package genomic_annotations

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
	"github.com/go-openapi/swag"
)

// NewGetValuesParams creates a new GetValuesParams object
// with the default values initialized.
func NewGetValuesParams() *GetValuesParams {
	var (
		limitDefault = int64(10)
	)
	return &GetValuesParams{
		Limit: &limitDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetValuesParamsWithTimeout creates a new GetValuesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetValuesParamsWithTimeout(timeout time.Duration) *GetValuesParams {
	var (
		limitDefault = int64(10)
	)
	return &GetValuesParams{
		Limit: &limitDefault,

		timeout: timeout,
	}
}

// NewGetValuesParamsWithContext creates a new GetValuesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetValuesParamsWithContext(ctx context.Context) *GetValuesParams {
	var (
		limitDefault = int64(10)
	)
	return &GetValuesParams{
		Limit: &limitDefault,

		Context: ctx,
	}
}

// NewGetValuesParamsWithHTTPClient creates a new GetValuesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetValuesParamsWithHTTPClient(client *http.Client) *GetValuesParams {
	var (
		limitDefault = int64(10)
	)
	return &GetValuesParams{
		Limit:      &limitDefault,
		HTTPClient: client,
	}
}

/*GetValuesParams contains all the parameters to send to the API endpoint
for the get values operation typically these are written to a http.Request
*/
type GetValuesParams struct {

	/*Annotation
	  Genomic annotation name.

	*/
	Annotation string
	/*Limit
	  Limits the number of records retrieved.

	*/
	Limit *int64
	/*Value
	  Genomic annotation value.

	*/
	Value string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get values params
func (o *GetValuesParams) WithTimeout(timeout time.Duration) *GetValuesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get values params
func (o *GetValuesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get values params
func (o *GetValuesParams) WithContext(ctx context.Context) *GetValuesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get values params
func (o *GetValuesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get values params
func (o *GetValuesParams) WithHTTPClient(client *http.Client) *GetValuesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get values params
func (o *GetValuesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAnnotation adds the annotation to the get values params
func (o *GetValuesParams) WithAnnotation(annotation string) *GetValuesParams {
	o.SetAnnotation(annotation)
	return o
}

// SetAnnotation adds the annotation to the get values params
func (o *GetValuesParams) SetAnnotation(annotation string) {
	o.Annotation = annotation
}

// WithLimit adds the limit to the get values params
func (o *GetValuesParams) WithLimit(limit *int64) *GetValuesParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get values params
func (o *GetValuesParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithValue adds the value to the get values params
func (o *GetValuesParams) WithValue(value string) *GetValuesParams {
	o.SetValue(value)
	return o
}

// SetValue adds the value to the get values params
func (o *GetValuesParams) SetValue(value string) {
	o.Value = value
}

// WriteToRequest writes these params to a swagger request
func (o *GetValuesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param annotation
	if err := r.SetPathParam("annotation", o.Annotation); err != nil {
		return err
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	// query param value
	qrValue := o.Value
	qValue := qrValue
	if qValue != "" {
		if err := r.SetQueryParam("value", qValue); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
