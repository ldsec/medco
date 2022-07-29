// Code generated by go-swagger; DO NOT EDIT.

package explore_statistics

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

// NewExploreStatisticsParams creates a new ExploreStatisticsParams object
// with the default values initialized.
func NewExploreStatisticsParams() *ExploreStatisticsParams {
	var ()
	return &ExploreStatisticsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewExploreStatisticsParamsWithTimeout creates a new ExploreStatisticsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewExploreStatisticsParamsWithTimeout(timeout time.Duration) *ExploreStatisticsParams {
	var ()
	return &ExploreStatisticsParams{

		timeout: timeout,
	}
}

// NewExploreStatisticsParamsWithContext creates a new ExploreStatisticsParams object
// with the default values initialized, and the ability to set a context for a request
func NewExploreStatisticsParamsWithContext(ctx context.Context) *ExploreStatisticsParams {
	var ()
	return &ExploreStatisticsParams{

		Context: ctx,
	}
}

// NewExploreStatisticsParamsWithHTTPClient creates a new ExploreStatisticsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewExploreStatisticsParamsWithHTTPClient(client *http.Client) *ExploreStatisticsParams {
	var ()
	return &ExploreStatisticsParams{
		HTTPClient: client,
	}
}

/*ExploreStatisticsParams contains all the parameters to send to the API endpoint
for the explore statistics operation typically these are written to a http.Request
*/
type ExploreStatisticsParams struct {

	/*Body
	  User public key, cohort name, modifier or concept information, interval size, minimum and maximum value of concept or modifier

	*/
	Body ExploreStatisticsBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the explore statistics params
func (o *ExploreStatisticsParams) WithTimeout(timeout time.Duration) *ExploreStatisticsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the explore statistics params
func (o *ExploreStatisticsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the explore statistics params
func (o *ExploreStatisticsParams) WithContext(ctx context.Context) *ExploreStatisticsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the explore statistics params
func (o *ExploreStatisticsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the explore statistics params
func (o *ExploreStatisticsParams) WithHTTPClient(client *http.Client) *ExploreStatisticsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the explore statistics params
func (o *ExploreStatisticsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the explore statistics params
func (o *ExploreStatisticsParams) WithBody(body ExploreStatisticsBody) *ExploreStatisticsParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the explore statistics params
func (o *ExploreStatisticsParams) SetBody(body ExploreStatisticsBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *ExploreStatisticsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}