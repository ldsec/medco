// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/ldsec/medco/connector/restapi/models"
)

// PostCohortsHandlerFunc turns a function with the right signature into a post cohorts handler
type PostCohortsHandlerFunc func(PostCohortsParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn PostCohortsHandlerFunc) Handle(params PostCohortsParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// PostCohortsHandler interface for that can handle valid post cohorts params
type PostCohortsHandler interface {
	Handle(PostCohortsParams, *models.User) middleware.Responder
}

// NewPostCohorts creates a new http.Handler for the post cohorts operation
func NewPostCohorts(ctx *middleware.Context, handler PostCohortsHandler) *PostCohorts {
	return &PostCohorts{Context: ctx, Handler: handler}
}

/*PostCohorts swagger:route POST /node/explore/cohorts/{name} medco-node postCohorts

Add a new cohort

*/
type PostCohorts struct {
	Context *middleware.Context
	Handler PostCohortsHandler
}

func (o *PostCohorts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostCohortsParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.User
	if uprinc != nil {
		principal = uprinc.(*models.User) // this is really a models.User, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostCohortsBadRequestBody post cohorts bad request body
//
// swagger:model PostCohortsBadRequestBody
type PostCohortsBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this post cohorts bad request body
func (o *PostCohortsBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCohortsBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCohortsBadRequestBody) UnmarshalBinary(b []byte) error {
	var res PostCohortsBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostCohortsBody post cohorts body
//
// swagger:model PostCohortsBody
type PostCohortsBody struct {

	// creation date
	// Required: true
	CreationDate *string `json:"creationDate"`

	// query ID
	// Required: true
	QueryID *int64 `json:"queryID"`

	// update date
	// Required: true
	UpdateDate *string `json:"updateDate"`
}

// Validate validates this post cohorts body
func (o *PostCohortsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCreationDate(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateQueryID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateUpdateDate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostCohortsBody) validateCreationDate(formats strfmt.Registry) error {

	if err := validate.Required("cohortsRequest"+"."+"creationDate", "body", o.CreationDate); err != nil {
		return err
	}

	return nil
}

func (o *PostCohortsBody) validateQueryID(formats strfmt.Registry) error {

	if err := validate.Required("cohortsRequest"+"."+"queryID", "body", o.QueryID); err != nil {
		return err
	}

	return nil
}

func (o *PostCohortsBody) validateUpdateDate(formats strfmt.Registry) error {

	if err := validate.Required("cohortsRequest"+"."+"updateDate", "body", o.UpdateDate); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostCohortsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCohortsBody) UnmarshalBinary(b []byte) error {
	var res PostCohortsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostCohortsConflictBody post cohorts conflict body
//
// swagger:model PostCohortsConflictBody
type PostCohortsConflictBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this post cohorts conflict body
func (o *PostCohortsConflictBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCohortsConflictBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCohortsConflictBody) UnmarshalBinary(b []byte) error {
	var res PostCohortsConflictBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostCohortsDefaultBody post cohorts default body
//
// swagger:model PostCohortsDefaultBody
type PostCohortsDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this post cohorts default body
func (o *PostCohortsDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCohortsDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCohortsDefaultBody) UnmarshalBinary(b []byte) error {
	var res PostCohortsDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostCohortsNotFoundBody post cohorts not found body
//
// swagger:model PostCohortsNotFoundBody
type PostCohortsNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this post cohorts not found body
func (o *PostCohortsNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostCohortsNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostCohortsNotFoundBody) UnmarshalBinary(b []byte) error {
	var res PostCohortsNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
