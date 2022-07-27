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

// GetDefaultCohortReader is a Reader for the GetDefaultCohort structure.
type GetDefaultCohortReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDefaultCohortReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDefaultCohortOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetDefaultCohortDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDefaultCohortOK creates a GetDefaultCohortOK with default headers values
func NewGetDefaultCohortOK() *GetDefaultCohortOK {
	return &GetDefaultCohortOK{}
}

/*GetDefaultCohortOK handles this case with default header values.

Name of cohort/filter retrieved
*/
type GetDefaultCohortOK struct {
	Payload string
}

func (o *GetDefaultCohortOK) Error() string {
	return fmt.Sprintf("[GET /node/explore/default-cohort][%d] getDefaultCohortOK  %+v", 200, o.Payload)
}

func (o *GetDefaultCohortOK) GetPayload() string {
	return o.Payload
}

func (o *GetDefaultCohortOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDefaultCohortDefault creates a GetDefaultCohortDefault with default headers values
func NewGetDefaultCohortDefault(code int) *GetDefaultCohortDefault {
	return &GetDefaultCohortDefault{
		_statusCode: code,
	}
}

/*GetDefaultCohortDefault handles this case with default header values.

Error response.
*/
type GetDefaultCohortDefault struct {
	_statusCode int

	Payload *GetDefaultCohortDefaultBody
}

// Code gets the status code for the get default cohort default response
func (o *GetDefaultCohortDefault) Code() int {
	return o._statusCode
}

func (o *GetDefaultCohortDefault) Error() string {
	return fmt.Sprintf("[GET /node/explore/default-cohort][%d] getDefaultCohort default  %+v", o._statusCode, o.Payload)
}

func (o *GetDefaultCohortDefault) GetPayload() *GetDefaultCohortDefaultBody {
	return o.Payload
}

func (o *GetDefaultCohortDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetDefaultCohortDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetDefaultCohortDefaultBody get default cohort default body
swagger:model GetDefaultCohortDefaultBody
*/
type GetDefaultCohortDefaultBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this get default cohort default body
func (o *GetDefaultCohortDefaultBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDefaultCohortDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDefaultCohortDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetDefaultCohortDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
