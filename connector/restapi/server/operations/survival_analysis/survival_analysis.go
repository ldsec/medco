// Code generated by go-swagger; DO NOT EDIT.

package survival_analysis

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/ldsec/medco/connector/restapi/models"
)

// SurvivalAnalysisHandlerFunc turns a function with the right signature into a survival analysis handler
type SurvivalAnalysisHandlerFunc func(SurvivalAnalysisParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn SurvivalAnalysisHandlerFunc) Handle(params SurvivalAnalysisParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// SurvivalAnalysisHandler interface for that can handle valid survival analysis params
type SurvivalAnalysisHandler interface {
	Handle(SurvivalAnalysisParams, *models.User) middleware.Responder
}

// NewSurvivalAnalysis creates a new http.Handler for the survival analysis operation
func NewSurvivalAnalysis(ctx *middleware.Context, handler SurvivalAnalysisHandler) *SurvivalAnalysis {
	return &SurvivalAnalysis{Context: ctx, Handler: handler}
}

/*SurvivalAnalysis swagger:route POST /node/analysis/survival/query survival-analysis survivalAnalysis

Send a query to run a survival analysis

*/
type SurvivalAnalysis struct {
	Context *middleware.Context
	Handler SurvivalAnalysisHandler
}

func (o *SurvivalAnalysis) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSurvivalAnalysisParams()

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

// SurvivalAnalysisBadRequestBody survival analysis bad request body
//
// swagger:model SurvivalAnalysisBadRequestBody
type SurvivalAnalysisBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this survival analysis bad request body
func (o *SurvivalAnalysisBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisBadRequestBody) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisBody survival analysis body
//
// swagger:model SurvivalAnalysisBody
type SurvivalAnalysisBody struct {

	// ID
	// Required: true
	// Pattern: ^[\w:-]+$
	ID *string `json:"ID"`

	// censoring from
	// Required: true
	// Enum: [encounters observations]
	CensoringFrom *string `json:"censoringFrom"`

	// cohort name
	// Required: true
	// Pattern: ^\w+$
	CohortName *string `json:"cohortName"`

	// end concept
	// Required: true
	// Pattern: ^\/$|^((\/[^\/]+)+\/?)$
	EndConcept *string `json:"endConcept"`

	// end modifier
	EndModifier *SurvivalAnalysisParamsBodyEndModifier `json:"endModifier,omitempty"`

	// ends when
	// Required: true
	// Enum: [earliest latest]
	EndsWhen *string `json:"endsWhen"`

	// start concept
	// Required: true
	// Pattern: ^\/$|^((\/[^\/]+)+\/?)$
	StartConcept *string `json:"startConcept"`

	// start modifier
	StartModifier *SurvivalAnalysisParamsBodyStartModifier `json:"startModifier,omitempty"`

	// starts when
	// Required: true
	// Enum: [earliest latest]
	StartsWhen *string `json:"startsWhen"`

	// sub group definitions
	// Max Items: 4
	SubGroupDefinitions []*SurvivalAnalysisParamsBodySubGroupDefinitionsItems0 `json:"subGroupDefinitions"`

	// time granularity
	// Required: true
	// Enum: [day week month year]
	TimeGranularity *string `json:"timeGranularity"`

	// time limit
	// Required: true
	// Minimum: 1
	TimeLimit *int64 `json:"timeLimit"`

	// user public key
	// Required: true
	// Pattern: ^[\w=-]+$
	UserPublicKey *string `json:"userPublicKey"`
}

// Validate validates this survival analysis body
func (o *SurvivalAnalysisBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCensoringFrom(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCohortName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateEndConcept(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateEndModifier(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateEndsWhen(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateStartConcept(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateStartModifier(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateStartsWhen(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSubGroupDefinitions(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimeGranularity(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimeLimit(formats); err != nil {
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

func (o *SurvivalAnalysisBody) validateID(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"ID", "body", o.ID); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"ID", "body", string(*o.ID), `^[\w:-]+$`); err != nil {
		return err
	}

	return nil
}

var survivalAnalysisBodyTypeCensoringFromPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["encounters","observations"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		survivalAnalysisBodyTypeCensoringFromPropEnum = append(survivalAnalysisBodyTypeCensoringFromPropEnum, v)
	}
}

const (

	// SurvivalAnalysisBodyCensoringFromEncounters captures enum value "encounters"
	SurvivalAnalysisBodyCensoringFromEncounters string = "encounters"

	// SurvivalAnalysisBodyCensoringFromObservations captures enum value "observations"
	SurvivalAnalysisBodyCensoringFromObservations string = "observations"
)

// prop value enum
func (o *SurvivalAnalysisBody) validateCensoringFromEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, survivalAnalysisBodyTypeCensoringFromPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SurvivalAnalysisBody) validateCensoringFrom(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"censoringFrom", "body", o.CensoringFrom); err != nil {
		return err
	}

	// value enum
	if err := o.validateCensoringFromEnum("body"+"."+"censoringFrom", "body", *o.CensoringFrom); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateCohortName(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"cohortName", "body", o.CohortName); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"cohortName", "body", string(*o.CohortName), `^\w+$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateEndConcept(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"endConcept", "body", o.EndConcept); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"endConcept", "body", string(*o.EndConcept), `^\/$|^((\/[^\/]+)+\/?)$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateEndModifier(formats strfmt.Registry) error {

	if swag.IsZero(o.EndModifier) { // not required
		return nil
	}

	if o.EndModifier != nil {
		if err := o.EndModifier.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "endModifier")
			}
			return err
		}
	}

	return nil
}

var survivalAnalysisBodyTypeEndsWhenPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["earliest","latest"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		survivalAnalysisBodyTypeEndsWhenPropEnum = append(survivalAnalysisBodyTypeEndsWhenPropEnum, v)
	}
}

const (

	// SurvivalAnalysisBodyEndsWhenEarliest captures enum value "earliest"
	SurvivalAnalysisBodyEndsWhenEarliest string = "earliest"

	// SurvivalAnalysisBodyEndsWhenLatest captures enum value "latest"
	SurvivalAnalysisBodyEndsWhenLatest string = "latest"
)

// prop value enum
func (o *SurvivalAnalysisBody) validateEndsWhenEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, survivalAnalysisBodyTypeEndsWhenPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SurvivalAnalysisBody) validateEndsWhen(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"endsWhen", "body", o.EndsWhen); err != nil {
		return err
	}

	// value enum
	if err := o.validateEndsWhenEnum("body"+"."+"endsWhen", "body", *o.EndsWhen); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateStartConcept(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"startConcept", "body", o.StartConcept); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"startConcept", "body", string(*o.StartConcept), `^\/$|^((\/[^\/]+)+\/?)$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateStartModifier(formats strfmt.Registry) error {

	if swag.IsZero(o.StartModifier) { // not required
		return nil
	}

	if o.StartModifier != nil {
		if err := o.StartModifier.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "startModifier")
			}
			return err
		}
	}

	return nil
}

var survivalAnalysisBodyTypeStartsWhenPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["earliest","latest"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		survivalAnalysisBodyTypeStartsWhenPropEnum = append(survivalAnalysisBodyTypeStartsWhenPropEnum, v)
	}
}

const (

	// SurvivalAnalysisBodyStartsWhenEarliest captures enum value "earliest"
	SurvivalAnalysisBodyStartsWhenEarliest string = "earliest"

	// SurvivalAnalysisBodyStartsWhenLatest captures enum value "latest"
	SurvivalAnalysisBodyStartsWhenLatest string = "latest"
)

// prop value enum
func (o *SurvivalAnalysisBody) validateStartsWhenEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, survivalAnalysisBodyTypeStartsWhenPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SurvivalAnalysisBody) validateStartsWhen(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"startsWhen", "body", o.StartsWhen); err != nil {
		return err
	}

	// value enum
	if err := o.validateStartsWhenEnum("body"+"."+"startsWhen", "body", *o.StartsWhen); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateSubGroupDefinitions(formats strfmt.Registry) error {

	if swag.IsZero(o.SubGroupDefinitions) { // not required
		return nil
	}

	iSubGroupDefinitionsSize := int64(len(o.SubGroupDefinitions))

	if err := validate.MaxItems("body"+"."+"subGroupDefinitions", "body", iSubGroupDefinitionsSize, 4); err != nil {
		return err
	}

	for i := 0; i < len(o.SubGroupDefinitions); i++ {
		if swag.IsZero(o.SubGroupDefinitions[i]) { // not required
			continue
		}

		if o.SubGroupDefinitions[i] != nil {
			if err := o.SubGroupDefinitions[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("body" + "." + "subGroupDefinitions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

var survivalAnalysisBodyTypeTimeGranularityPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["day","week","month","year"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		survivalAnalysisBodyTypeTimeGranularityPropEnum = append(survivalAnalysisBodyTypeTimeGranularityPropEnum, v)
	}
}

const (

	// SurvivalAnalysisBodyTimeGranularityDay captures enum value "day"
	SurvivalAnalysisBodyTimeGranularityDay string = "day"

	// SurvivalAnalysisBodyTimeGranularityWeek captures enum value "week"
	SurvivalAnalysisBodyTimeGranularityWeek string = "week"

	// SurvivalAnalysisBodyTimeGranularityMonth captures enum value "month"
	SurvivalAnalysisBodyTimeGranularityMonth string = "month"

	// SurvivalAnalysisBodyTimeGranularityYear captures enum value "year"
	SurvivalAnalysisBodyTimeGranularityYear string = "year"
)

// prop value enum
func (o *SurvivalAnalysisBody) validateTimeGranularityEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, survivalAnalysisBodyTypeTimeGranularityPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SurvivalAnalysisBody) validateTimeGranularity(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"timeGranularity", "body", o.TimeGranularity); err != nil {
		return err
	}

	// value enum
	if err := o.validateTimeGranularityEnum("body"+"."+"timeGranularity", "body", *o.TimeGranularity); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateTimeLimit(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"timeLimit", "body", o.TimeLimit); err != nil {
		return err
	}

	if err := validate.MinimumInt("body"+"."+"timeLimit", "body", int64(*o.TimeLimit), 1, false); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisBody) validateUserPublicKey(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"userPublicKey", "body", o.UserPublicKey); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"userPublicKey", "body", string(*o.UserPublicKey), `^[\w=-]+$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisBody) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisDefaultBody survival analysis default body
//
// swagger:model SurvivalAnalysisDefaultBody
type SurvivalAnalysisDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this survival analysis default body
func (o *SurvivalAnalysisDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisDefaultBody) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisNotFoundBody survival analysis not found body
//
// swagger:model SurvivalAnalysisNotFoundBody
type SurvivalAnalysisNotFoundBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this survival analysis not found body
func (o *SurvivalAnalysisNotFoundBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisNotFoundBody) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisOKBody survival analysis o k body
//
// swagger:model SurvivalAnalysisOKBody
type SurvivalAnalysisOKBody struct {

	// results
	Results []*SurvivalAnalysisOKBodyResultsItems0 `json:"results"`

	// timers
	Timers models.Timers `json:"timers"`
}

// Validate validates this survival analysis o k body
func (o *SurvivalAnalysisOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResults(formats); err != nil {
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

func (o *SurvivalAnalysisOKBody) validateResults(formats strfmt.Registry) error {

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
					return ve.ValidateName("survivalAnalysisOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *SurvivalAnalysisOKBody) validateTimers(formats strfmt.Registry) error {

	if swag.IsZero(o.Timers) { // not required
		return nil
	}

	if err := o.Timers.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("survivalAnalysisOK" + "." + "timers")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisOKBody) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisOKBodyResultsItems0 survival analysis o k body results items0
//
// swagger:model SurvivalAnalysisOKBodyResultsItems0
type SurvivalAnalysisOKBodyResultsItems0 struct {

	// group ID
	GroupID string `json:"groupID,omitempty"`

	// group results
	GroupResults []*SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0 `json:"groupResults"`

	// initial count
	InitialCount string `json:"initialCount,omitempty"`
}

// Validate validates this survival analysis o k body results items0
func (o *SurvivalAnalysisOKBodyResultsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGroupResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SurvivalAnalysisOKBodyResultsItems0) validateGroupResults(formats strfmt.Registry) error {

	if swag.IsZero(o.GroupResults) { // not required
		return nil
	}

	for i := 0; i < len(o.GroupResults); i++ {
		if swag.IsZero(o.GroupResults[i]) { // not required
			continue
		}

		if o.GroupResults[i] != nil {
			if err := o.GroupResults[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("groupResults" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisOKBodyResultsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0 survival analysis o k body results items0 group results items0
//
// swagger:model SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0
type SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0 struct {

	// events
	// Required: true
	Events *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events `json:"events"`

	// timepoint
	// Required: true
	Timepoint *int64 `json:"timepoint"`
}

// Validate validates this survival analysis o k body results items0 group results items0
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEvents(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimepoint(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0) validateEvents(formats strfmt.Registry) error {

	if err := validate.Required("events", "body", o.Events); err != nil {
		return err
	}

	if o.Events != nil {
		if err := o.Events.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("events")
			}
			return err
		}
	}

	return nil
}

func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0) validateTimepoint(formats strfmt.Registry) error {

	if err := validate.Required("timepoint", "body", o.Timepoint); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events survival analysis o k body results items0 group results items0 events
//
// swagger:model SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events
type SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events struct {

	// censoringevent
	// Required: true
	Censoringevent *string `json:"censoringevent"`

	// eventofinterest
	// Required: true
	Eventofinterest *string `json:"eventofinterest"`
}

// Validate validates this survival analysis o k body results items0 group results items0 events
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCensoringevent(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateEventofinterest(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events) validateCensoringevent(formats strfmt.Registry) error {

	if err := validate.Required("events"+"."+"censoringevent", "body", o.Censoringevent); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events) validateEventofinterest(formats strfmt.Registry) error {

	if err := validate.Required("events"+"."+"eventofinterest", "body", o.Eventofinterest); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisParamsBodyEndModifier survival analysis params body end modifier
//
// swagger:model SurvivalAnalysisParamsBodyEndModifier
type SurvivalAnalysisParamsBodyEndModifier struct {

	// applied path
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/%?)$
	AppliedPath *string `json:"appliedPath"`

	// modifier key
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/)$
	ModifierKey *string `json:"modifierKey"`
}

// Validate validates this survival analysis params body end modifier
func (o *SurvivalAnalysisParamsBodyEndModifier) Validate(formats strfmt.Registry) error {
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

func (o *SurvivalAnalysisParamsBodyEndModifier) validateAppliedPath(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"endModifier"+"."+"appliedPath", "body", o.AppliedPath); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"endModifier"+"."+"appliedPath", "body", string(*o.AppliedPath), `^((\/[^\/]+)+\/%?)$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisParamsBodyEndModifier) validateModifierKey(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"endModifier"+"."+"modifierKey", "body", o.ModifierKey); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"endModifier"+"."+"modifierKey", "body", string(*o.ModifierKey), `^((\/[^\/]+)+\/)$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodyEndModifier) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodyEndModifier) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisParamsBodyEndModifier
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisParamsBodyStartModifier survival analysis params body start modifier
//
// swagger:model SurvivalAnalysisParamsBodyStartModifier
type SurvivalAnalysisParamsBodyStartModifier struct {

	// applied path
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/%?)$
	AppliedPath *string `json:"appliedPath"`

	// modifier key
	// Required: true
	// Pattern: ^((\/[^\/]+)+\/)$
	ModifierKey *string `json:"modifierKey"`
}

// Validate validates this survival analysis params body start modifier
func (o *SurvivalAnalysisParamsBodyStartModifier) Validate(formats strfmt.Registry) error {
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

func (o *SurvivalAnalysisParamsBodyStartModifier) validateAppliedPath(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"startModifier"+"."+"appliedPath", "body", o.AppliedPath); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"startModifier"+"."+"appliedPath", "body", string(*o.AppliedPath), `^((\/[^\/]+)+\/%?)$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisParamsBodyStartModifier) validateModifierKey(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"startModifier"+"."+"modifierKey", "body", o.ModifierKey); err != nil {
		return err
	}

	if err := validate.Pattern("body"+"."+"startModifier"+"."+"modifierKey", "body", string(*o.ModifierKey), `^((\/[^\/]+)+\/)$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodyStartModifier) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodyStartModifier) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisParamsBodyStartModifier
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SurvivalAnalysisParamsBodySubGroupDefinitionsItems0 survival analysis params body sub group definitions items0
//
// swagger:model SurvivalAnalysisParamsBodySubGroupDefinitionsItems0
type SurvivalAnalysisParamsBodySubGroupDefinitionsItems0 struct {

	// group name
	// Pattern: ^\w+$
	GroupName string `json:"groupName,omitempty"`

	// panels
	Panels []*models.Panel `json:"panels"`

	// sub group timing
	SubGroupTiming models.Timing `json:"subGroupTiming,omitempty"`
}

// Validate validates this survival analysis params body sub group definitions items0
func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGroupName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePanels(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSubGroupTiming(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) validateGroupName(formats strfmt.Registry) error {

	if swag.IsZero(o.GroupName) { // not required
		return nil
	}

	if err := validate.Pattern("groupName", "body", string(o.GroupName), `^\w+$`); err != nil {
		return err
	}

	return nil
}

func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) validatePanels(formats strfmt.Registry) error {

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
					return ve.ValidateName("panels" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) validateSubGroupTiming(formats strfmt.Registry) error {

	if swag.IsZero(o.SubGroupTiming) { // not required
		return nil
	}

	if err := o.SubGroupTiming.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("subGroupTiming")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) UnmarshalBinary(b []byte) error {
	var res SurvivalAnalysisParamsBodySubGroupDefinitionsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
