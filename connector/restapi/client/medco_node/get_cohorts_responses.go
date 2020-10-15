// Code generated by go-swagger; DO NOT EDIT.

package medco_node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetCohortsReader is a Reader for the GetCohorts structure.
type GetCohortsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCohortsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCohortsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetCohortsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetCohortsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCohortsOK creates a GetCohortsOK with default headers values
func NewGetCohortsOK() *GetCohortsOK {
	return &GetCohortsOK{}
}

/*GetCohortsOK handles this case with default header values.

Queried cohorts
*/
type GetCohortsOK struct {
	Payload []*GetCohortsOKBodyItems0
}

func (o *GetCohortsOK) Error() string {
	return fmt.Sprintf("[GET /node/explore/cohorts][%d] getCohortsOK  %+v", 200, o.Payload)
}

func (o *GetCohortsOK) GetPayload() []*GetCohortsOKBodyItems0 {
	return o.Payload
}

func (o *GetCohortsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCohortsNotFound creates a GetCohortsNotFound with default headers values
func NewGetCohortsNotFound() *GetCohortsNotFound {
	return &GetCohortsNotFound{}
}

/*GetCohortsNotFound handles this case with default header values.

User not found
*/
type GetCohortsNotFound struct {
}

func (o *GetCohortsNotFound) Error() string {
	return fmt.Sprintf("[GET /node/explore/cohorts][%d] getCohortsNotFound ", 404)
}

func (o *GetCohortsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCohortsDefault creates a GetCohortsDefault with default headers values
func NewGetCohortsDefault(code int) *GetCohortsDefault {
	return &GetCohortsDefault{
		_statusCode: code,
	}
}

/*GetCohortsDefault handles this case with default header values.

Error response.
*/
type GetCohortsDefault struct {
	_statusCode int

	Payload *GetCohortsDefaultBody
}

// Code gets the status code for the get cohorts default response
func (o *GetCohortsDefault) Code() int {
	return o._statusCode
}

func (o *GetCohortsDefault) Error() string {
	return fmt.Sprintf("[GET /node/explore/cohorts][%d] getCohorts default  %+v", o._statusCode, o.Payload)
}

func (o *GetCohortsDefault) GetPayload() *GetCohortsDefaultBody {
	return o.Payload
}

func (o *GetCohortsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetCohortsDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetCohortsDefaultBody get cohorts default body
swagger:model GetCohortsDefaultBody
*/
type GetCohortsDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get cohorts default body
func (o *GetCohortsDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetCohortsDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetCohortsDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetCohortsDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetCohortsOKBodyItems0 get cohorts o k body items0
swagger:model GetCohortsOKBodyItems0
*/
type GetCohortsOKBodyItems0 struct {

	// cohort Id
	CohortID int64 `json:"cohortId,omitempty"`

	// cohort name
	CohortName string `json:"cohortName,omitempty"`

	// creation date
	CreationDate string `json:"creationDate,omitempty"`

	// query Id
	QueryID int64 `json:"queryId,omitempty"`

	// update date
	UpdateDate string `json:"updateDate,omitempty"`
}

// Validate validates this get cohorts o k body items0
func (o *GetCohortsOKBodyItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetCohortsOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetCohortsOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res GetCohortsOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}