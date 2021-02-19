// Code generated by go-swagger; DO NOT EDIT.

package explore_statistics

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/ldsec/medco/connector/restapi/models"
)

// ExploreStatisticsHandlerFunc turns a function with the right signature into a explore statistics handler
type ExploreStatisticsHandlerFunc func(ExploreStatisticsParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn ExploreStatisticsHandlerFunc) Handle(params ExploreStatisticsParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// ExploreStatisticsHandler interface for that can handle valid explore statistics params
type ExploreStatisticsHandler interface {
	Handle(ExploreStatisticsParams, *models.User) middleware.Responder
}

// NewExploreStatistics creates a new http.Handler for the explore statistics operation
func NewExploreStatistics(ctx *middleware.Context, handler ExploreStatisticsHandler) *ExploreStatistics {
	return &ExploreStatistics{Context: ctx, Handler: handler}
}

/*ExploreStatistics swagger:route POST /node/explore-statistics/query explore statistics exploreStatistics

Queries the server to obtain the histogram of the distribution linked to the concept sent as parameter to the request.

*/
type ExploreStatistics struct {
	Context *middleware.Context
	Handler ExploreStatisticsHandler
}

func (o *ExploreStatistics) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewExploreStatisticsParams()

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

// ExploreStatisticsBadRequestBody explore statistics bad request body
//
// swagger:model ExploreStatisticsBadRequestBody
type ExploreStatisticsBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this explore statistics bad request body
func (o *ExploreStatisticsBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsBadRequestBody) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsBody explore statistics body
//
// swagger:model ExploreStatisticsBody
type ExploreStatisticsBody struct {

	// ID
	// Pattern: ^[\w:-]+$
	ID string `json:"ID,omitempty"`

	// bucket size
	BucketSize float64 `json:"bucketSize,omitempty"`

	// cohort definition
	CohortDefinition *ExploreStatisticsParamsBodyCohortDefinition `json:"cohortDefinition,omitempty"`

	// A list of the paths of concepts used as analytes
	Concepts []string `json:"concepts"`

	// min observation
	MinObservation float64 `json:"minObservation,omitempty"`

	// A list describing the modifiers used as analytes
	Modifiers []*ExploreStatisticsParamsBodyModifiersItems0 `json:"modifiers"`

	// user public key
	// Pattern: ^[\w=-]+$
	UserPublicKey string `json:"userPublicKey,omitempty"`
}

// Validate validates this explore statistics body
func (o *ExploreStatisticsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCohortDefinition(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateConcepts(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateModifiers(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateUserPublicKey(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ExploreStatisticsBody) validateID(formats strfmt.Registry) error {

	if swag.IsZero(o.ID) { // not required
		return nil
	}

	if err := validate.Pattern("body"+"."+"ID", "body", string(o.ID), `^[\w:-]+$`); err != nil {
		return err
	}

	return nil
}

func (o *ExploreStatisticsBody) validateCohortDefinition(formats strfmt.Registry) error {

	if swag.IsZero(o.CohortDefinition) { // not required
		return nil
	}

	if o.CohortDefinition != nil {
		if err := o.CohortDefinition.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "cohortDefinition")
			}
			return err
		}
	}

	return nil
}

func (o *ExploreStatisticsBody) validateConcepts(formats strfmt.Registry) error {

	if swag.IsZero(o.Concepts) { // not required
		return nil
	}

	for i := 0; i < len(o.Concepts); i++ {

		if err := validate.Pattern("body"+"."+"concepts"+"."+strconv.Itoa(i), "body", string(o.Concepts[i]), `^\/$|^((\/[^\/]+)+\/?)$`); err != nil {
			return err
		}

	}

	return nil
}

func (o *ExploreStatisticsBody) validateModifiers(formats strfmt.Registry) error {

	if swag.IsZero(o.Modifiers) { // not required
		return nil
	}

	for i := 0; i < len(o.Modifiers); i++ {
		if swag.IsZero(o.Modifiers[i]) { // not required
			continue
		}

		if o.Modifiers[i] != nil {
			if err := o.Modifiers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("body" + "." + "modifiers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ExploreStatisticsBody) validateUserPublicKey(formats strfmt.Registry) error {

	if swag.IsZero(o.UserPublicKey) { // not required
		return nil
	}

	if err := validate.Pattern("body"+"."+"userPublicKey", "body", string(o.UserPublicKey), `^[\w=-]+$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsBody) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsDefaultBody explore statistics default body
//
// swagger:model ExploreStatisticsDefaultBody
type ExploreStatisticsDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this explore statistics default body
func (o *ExploreStatisticsDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsDefaultBody) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsNotFoundBody explore statistics not found body
//
// swagger:model ExploreStatisticsNotFoundBody
type ExploreStatisticsNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this explore statistics not found body
func (o *ExploreStatisticsNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsNotFoundBody) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsOKBody explore statistics o k body
//
// swagger:model ExploreStatisticsOKBody
type ExploreStatisticsOKBody struct {

	// Timers for work happening outside of the construction of the histograms
	GlobalTimers models.Timers `json:"globalTimers"`

	// Each item of this array contains the histogram of a specific analyte (concept or modifier).
	Results []*ExploreStatisticsOKBodyResultsItems0 `json:"results"`
}

// Validate validates this explore statistics o k body
func (o *ExploreStatisticsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGlobalTimers(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ExploreStatisticsOKBody) validateGlobalTimers(formats strfmt.Registry) error {

	if swag.IsZero(o.GlobalTimers) { // not required
		return nil
	}

	if err := o.GlobalTimers.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("exploreStatisticsOK" + "." + "globalTimers")
		}
		return err
	}

	return nil
}

func (o *ExploreStatisticsOKBody) validateResults(formats strfmt.Registry) error {

	if swag.IsZero(o.Results) { // not required
		return nil
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("exploreStatisticsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsOKBody) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsOKBodyResultsItems0 explore statistics o k body results items0
//
// swagger:model ExploreStatisticsOKBodyResultsItems0
type ExploreStatisticsOKBodyResultsItems0 struct {

	// The name of the analyte used to build this histogram
	AnalyteName string `json:"analyteName,omitempty"`

	// the encrypted counts of each bucket of the histogram
	Intervals []*models.IntervalBucket `json:"intervals"`

	// timers
	Timers models.Timers `json:"timers,omitempty"`

	// unit
	Unit string `json:"unit,omitempty"`
}

// Validate validates this explore statistics o k body results items0
func (o *ExploreStatisticsOKBodyResultsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateIntervals(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ExploreStatisticsOKBodyResultsItems0) validateIntervals(formats strfmt.Registry) error {

	if swag.IsZero(o.Intervals) { // not required
		return nil
	}

	for i := 0; i < len(o.Intervals); i++ {
		if swag.IsZero(o.Intervals[i]) { // not required
			continue
		}

		if o.Intervals[i] != nil {
			if err := o.Intervals[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("intervals" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ExploreStatisticsOKBodyResultsItems0) validateTimers(formats strfmt.Registry) error {

	if swag.IsZero(o.Timers) { // not required
		return nil
	}

	if err := o.Timers.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("timers")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsOKBodyResultsItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsOKBodyResultsItems0) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsOKBodyResultsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsParamsBodyCohortDefinition explore statistics params body cohort definition
//
// swagger:model ExploreStatisticsParamsBodyCohortDefinition
type ExploreStatisticsParamsBodyCohortDefinition struct {

	// i2b2 panels (linked by an AND)
	Panels []*models.Panel `json:"panels"`

	// query timing
	QueryTiming models.Timing `json:"queryTiming,omitempty"`
}

// Validate validates this explore statistics params body cohort definition
func (o *ExploreStatisticsParamsBodyCohortDefinition) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePanels(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateQueryTiming(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ExploreStatisticsParamsBodyCohortDefinition) validatePanels(formats strfmt.Registry) error {

	if swag.IsZero(o.Panels) { // not required
		return nil
	}

	for i := 0; i < len(o.Panels); i++ {
		if swag.IsZero(o.Panels[i]) { // not required
			continue
		}

		if o.Panels[i] != nil {
			if err := o.Panels[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("body" + "." + "cohortDefinition" + "." + "panels" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ExploreStatisticsParamsBodyCohortDefinition) validateQueryTiming(formats strfmt.Registry) error {

	if swag.IsZero(o.QueryTiming) { // not required
		return nil
	}

	if err := o.QueryTiming.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body" + "." + "cohortDefinition" + "." + "queryTiming")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsParamsBodyCohortDefinition) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsParamsBodyCohortDefinition) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsParamsBodyCohortDefinition
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// ExploreStatisticsParamsBodyModifiersItems0 explore statistics params body modifiers items0
//
// swagger:model ExploreStatisticsParamsBodyModifiersItems0
type ExploreStatisticsParamsBodyModifiersItems0 struct {

	// applied path
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/%?)$
	AppliedPath *string `json:"appliedPath"`

	// modifier key
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/)$
	ModifierKey *string `json:"modifierKey"`
}

// Validate validates this explore statistics params body modifiers items0
func (o *ExploreStatisticsParamsBodyModifiersItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAppliedPath(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateModifierKey(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ExploreStatisticsParamsBodyModifiersItems0) validateAppliedPath(formats strfmt.Registry) error {

	if err := validate.Required("appliedPath", "body", o.AppliedPath); err != nil {
		return err
	}

	if err := validate.Pattern("appliedPath", "body", string(*o.AppliedPath), `^((\/[^\/]+)+\/%?)$`); err != nil {
		return err
	}

	return nil
}

func (o *ExploreStatisticsParamsBodyModifiersItems0) validateModifierKey(formats strfmt.Registry) error {

	if err := validate.Required("modifierKey", "body", o.ModifierKey); err != nil {
		return err
	}

	if err := validate.Pattern("modifierKey", "body", string(*o.ModifierKey), `^((\/[^\/]+)+\/)$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ExploreStatisticsParamsBodyModifiersItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExploreStatisticsParamsBodyModifiersItems0) UnmarshalBinary(b []byte) error {
	var res ExploreStatisticsParamsBodyModifiersItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
