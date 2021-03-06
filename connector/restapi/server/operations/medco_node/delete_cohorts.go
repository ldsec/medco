// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/ldsec/medco/connector/restapi/models"
)

// DeleteCohortsHandlerFunc turns a function with the right signature into a delete cohorts handler
type DeleteCohortsHandlerFunc func(DeleteCohortsParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteCohortsHandlerFunc) Handle(params DeleteCohortsParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// DeleteCohortsHandler interface for that can handle valid delete cohorts params
type DeleteCohortsHandler interface {
	Handle(DeleteCohortsParams, *models.User) middleware.Responder
}

// NewDeleteCohorts creates a new http.Handler for the delete cohorts operation
func NewDeleteCohorts(ctx *middleware.Context, handler DeleteCohortsHandler) *DeleteCohorts {
	return &DeleteCohorts{Context: ctx, Handler: handler}
}

/*DeleteCohorts swagger:route DELETE /node/explore/cohorts/{name} medco-node deleteCohorts

Delete a cohort if it exists

*/
type DeleteCohorts struct {
	Context *middleware.Context
	Handler DeleteCohortsHandler
}

func (o *DeleteCohorts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteCohortsParams()

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

// DeleteCohortsDefaultBody delete cohorts default body
//
// swagger:model DeleteCohortsDefaultBody
type DeleteCohortsDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this delete cohorts default body
func (o *DeleteCohortsDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DeleteCohortsDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteCohortsDefaultBody) UnmarshalBinary(b []byte) error {
	var res DeleteCohortsDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// DeleteCohortsNotFoundBody delete cohorts not found body
//
// swagger:model DeleteCohortsNotFoundBody
type DeleteCohortsNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this delete cohorts not found body
func (o *DeleteCohortsNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *DeleteCohortsNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeleteCohortsNotFoundBody) UnmarshalBinary(b []byte) error {
	var res DeleteCohortsNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
