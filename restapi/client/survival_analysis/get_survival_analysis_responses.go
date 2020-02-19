// Code generated by go-swagger; DO NOT EDIT.

package survival_analysis

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/ldsec/medco-connector/restapi/models"
)

// GetSurvivalAnalysisReader is a Reader for the GetSurvivalAnalysis structure.
type GetSurvivalAnalysisReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSurvivalAnalysisReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetSurvivalAnalysisOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetSurvivalAnalysisNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetSurvivalAnalysisDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetSurvivalAnalysisOK creates a GetSurvivalAnalysisOK with default headers values
func NewGetSurvivalAnalysisOK() *GetSurvivalAnalysisOK {
	return &GetSurvivalAnalysisOK{}
}

/*GetSurvivalAnalysisOK handles this case with default header values.

Queried survival analysis
*/
type GetSurvivalAnalysisOK struct {
	Payload []*GetSurvivalAnalysisOKBodyItems0
}

func (o *GetSurvivalAnalysisOK) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{granularity}][%d] getSurvivalAnalysisOK  %+v", 200, o.Payload)
}

func (o *GetSurvivalAnalysisOK) GetPayload() []*GetSurvivalAnalysisOKBodyItems0 {
	return o.Payload
}

func (o *GetSurvivalAnalysisOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSurvivalAnalysisNotFound creates a GetSurvivalAnalysisNotFound with default headers values
func NewGetSurvivalAnalysisNotFound() *GetSurvivalAnalysisNotFound {
	return &GetSurvivalAnalysisNotFound{}
}

/*GetSurvivalAnalysisNotFound handles this case with default header values.

TODO not found
*/
type GetSurvivalAnalysisNotFound struct {
}

func (o *GetSurvivalAnalysisNotFound) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{granularity}][%d] getSurvivalAnalysisNotFound ", 404)
}

func (o *GetSurvivalAnalysisNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetSurvivalAnalysisDefault creates a GetSurvivalAnalysisDefault with default headers values
func NewGetSurvivalAnalysisDefault(code int) *GetSurvivalAnalysisDefault {
	return &GetSurvivalAnalysisDefault{
		_statusCode: code,
	}
}

/*GetSurvivalAnalysisDefault handles this case with default header values.

Error response.
*/
type GetSurvivalAnalysisDefault struct {
	_statusCode int

	Payload *GetSurvivalAnalysisDefaultBody
}

// Code gets the status code for the get survival analysis default response
func (o *GetSurvivalAnalysisDefault) Code() int {
	return o._statusCode
}

func (o *GetSurvivalAnalysisDefault) Error() string {
	return fmt.Sprintf("[GET /survival-analysis/{granularity}][%d] getSurvivalAnalysis default  %+v", o._statusCode, o.Payload)
}

func (o *GetSurvivalAnalysisDefault) GetPayload() *GetSurvivalAnalysisDefaultBody {
	return o.Payload
}

func (o *GetSurvivalAnalysisDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetSurvivalAnalysisDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetSurvivalAnalysisBody get survival analysis body
swagger:model GetSurvivalAnalysisBody
*/
type GetSurvivalAnalysisBody struct {

	// panels
	Panels models.ExploreQueryPanels `json:"panels,omitempty"`

	// user public key
	// Pattern: ^[\w=-]+$
	UserPublicKey string `json:"userPublicKey,omitempty"`
}

// Validate validates this get survival analysis body
func (o *GetSurvivalAnalysisBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePanels(formats); err != nil {
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

func (o *GetSurvivalAnalysisBody) validatePanels(formats strfmt.Registry) error {

	if swag.IsZero(o.Panels) { // not required
		return nil
	}

	if err := o.Panels.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("userPublicKeyAndPanels" + "." + "panels")
		}
		return err
	}

	return nil
}

func (o *GetSurvivalAnalysisBody) validateUserPublicKey(formats strfmt.Registry) error {

	if swag.IsZero(o.UserPublicKey) { // not required
		return nil
	}

	if err := validate.Pattern("userPublicKeyAndPanels"+"."+"userPublicKey", "body", string(o.UserPublicKey), `^[\w=-]+$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetSurvivalAnalysisBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetSurvivalAnalysisBody) UnmarshalBinary(b []byte) error {
	var res GetSurvivalAnalysisBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetSurvivalAnalysisDefaultBody get survival analysis default body
swagger:model GetSurvivalAnalysisDefaultBody
*/
type GetSurvivalAnalysisDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get survival analysis default body
func (o *GetSurvivalAnalysisDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetSurvivalAnalysisDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetSurvivalAnalysisDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetSurvivalAnalysisDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetSurvivalAnalysisOKBodyItems0 get survival analysis o k body items0
swagger:model GetSurvivalAnalysisOKBodyItems0
*/
type GetSurvivalAnalysisOKBodyItems0 struct {

	// events
	Events *GetSurvivalAnalysisOKBodyItems0Events `json:"events,omitempty"`

	// timepoint
	Timepoint string `json:"timepoint,omitempty"`
}

// Validate validates this get survival analysis o k body items0
func (o *GetSurvivalAnalysisOKBodyItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEvents(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetSurvivalAnalysisOKBodyItems0) validateEvents(formats strfmt.Registry) error {

	if swag.IsZero(o.Events) { // not required
		return nil
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

// MarshalBinary interface implementation
func (o *GetSurvivalAnalysisOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetSurvivalAnalysisOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res GetSurvivalAnalysisOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetSurvivalAnalysisOKBodyItems0Events get survival analysis o k body items0 events
swagger:model GetSurvivalAnalysisOKBodyItems0Events
*/
type GetSurvivalAnalysisOKBodyItems0Events struct {

	// censoringevent
	Censoringevent string `json:"censoringevent,omitempty"`

	// eventofinterest
	Eventofinterest string `json:"eventofinterest,omitempty"`
}

// Validate validates this get survival analysis o k body items0 events
func (o *GetSurvivalAnalysisOKBodyItems0Events) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetSurvivalAnalysisOKBodyItems0Events) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetSurvivalAnalysisOKBodyItems0Events) UnmarshalBinary(b []byte) error {
	var res GetSurvivalAnalysisOKBodyItems0Events
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
