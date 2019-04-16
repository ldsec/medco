// Code generated by go-swagger; DO NOT EDIT.

package picsure2

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	models "github.com/lca1/medco-connector/swagger/models"
)

// QuerySyncHandlerFunc turns a function with the right signature into a query sync handler
type QuerySyncHandlerFunc func(QuerySyncParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn QuerySyncHandlerFunc) Handle(params QuerySyncParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// QuerySyncHandler interface for that can handle valid query sync params
type QuerySyncHandler interface {
	Handle(QuerySyncParams, *models.User) middleware.Responder
}

// NewQuerySync creates a new http.Handler for the query sync operation
func NewQuerySync(ctx *middleware.Context, handler QuerySyncHandler) *QuerySync {
	return &QuerySync{Context: ctx, Handler: handler}
}

/*QuerySync swagger:route POST /picsure2/query/sync picsure2 querySync

Query MedCo node synchronously.

*/
type QuerySync struct {
	Context *middleware.Context
	Handler QuerySyncHandler
}

func (o *QuerySync) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewQuerySyncParams()

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

// QuerySyncBody query sync body
// swagger:model QuerySyncBody
type QuerySyncBody struct {

	// query
	Query *models.Query `json:"query,omitempty"`

	// resource credentials
	ResourceCredentials *models.ResourceCredentials `json:"resourceCredentials,omitempty"`

	// resource UUID
	ResourceUUID string `json:"resourceUUID,omitempty"`
}

// Validate validates this query sync body
func (o *QuerySyncBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateQuery(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResourceCredentials(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *QuerySyncBody) validateQuery(formats strfmt.Registry) error {

	if swag.IsZero(o.Query) { // not required
		return nil
	}

	if o.Query != nil {
		if err := o.Query.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "query")
			}
			return err
		}
	}

	return nil
}

func (o *QuerySyncBody) validateResourceCredentials(formats strfmt.Registry) error {

	if swag.IsZero(o.ResourceCredentials) { // not required
		return nil
	}

	if o.ResourceCredentials != nil {
		if err := o.ResourceCredentials.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "resourceCredentials")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *QuerySyncBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *QuerySyncBody) UnmarshalBinary(b []byte) error {
	var res QuerySyncBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// QuerySyncDefaultBody query sync default body
// swagger:model QuerySyncDefaultBody
type QuerySyncDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this query sync default body
func (o *QuerySyncDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *QuerySyncDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *QuerySyncDefaultBody) UnmarshalBinary(b []byte) error {
	var res QuerySyncDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
