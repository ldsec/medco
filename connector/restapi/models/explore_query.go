// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ExploreQuery MedCo-Explore query
//
// swagger:model exploreQuery
type ExploreQuery struct {

	// i2b2 panels (linked by an AND)
	Panels []*ExploreQueryPanelsItems0 `json:"panels"`

	// type
	Type ExploreQueryType `json:"type,omitempty"`

	// user public key
	// Pattern: ^[\w=-]+$
	UserPublicKey string `json:"userPublicKey,omitempty"`
}

// Validate validates this explore query
func (m *ExploreQuery) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePanels(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserPublicKey(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExploreQuery) validatePanels(formats strfmt.Registry) error {

	if swag.IsZero(m.Panels) { // not required
		return nil
	}

	for i := 0; i < len(m.Panels); i++ {
		if swag.IsZero(m.Panels[i]) { // not required
			continue
		}

		if m.Panels[i] != nil {
			if err := m.Panels[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("panels" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ExploreQuery) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if err := m.Type.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("type")
		}
		return err
	}

	return nil
}

func (m *ExploreQuery) validateUserPublicKey(formats strfmt.Registry) error {

	if swag.IsZero(m.UserPublicKey) { // not required
		return nil
	}

	if err := validate.Pattern("userPublicKey", "body", string(m.UserPublicKey), `^[\w=-]+$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExploreQuery) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExploreQuery) UnmarshalBinary(b []byte) error {
	var res ExploreQuery
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ExploreQueryPanelsItems0 explore query panels items0
//
// swagger:model ExploreQueryPanelsItems0
type ExploreQueryPanelsItems0 struct {

	// i2b2 items (linked by an OR)
	Items []*ExploreQueryPanelsItems0ItemsItems0 `json:"items"`

	// exclude the i2b2 panel
	// Required: true
	Not *bool `json:"not"`
}

// Validate validates this explore query panels items0
func (m *ExploreQueryPanelsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNot(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExploreQueryPanelsItems0) validateItems(formats strfmt.Registry) error {

	if swag.IsZero(m.Items) { // not required
		return nil
	}

	for i := 0; i < len(m.Items); i++ {
		if swag.IsZero(m.Items[i]) { // not required
			continue
		}

		if m.Items[i] != nil {
			if err := m.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ExploreQueryPanelsItems0) validateNot(formats strfmt.Registry) error {

	if err := validate.Required("not", "body", m.Not); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExploreQueryPanelsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExploreQueryPanelsItems0) UnmarshalBinary(b []byte) error {
	var res ExploreQueryPanelsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ExploreQueryPanelsItems0ItemsItems0 explore query panels items0 items items0
//
// swagger:model ExploreQueryPanelsItems0ItemsItems0
type ExploreQueryPanelsItems0ItemsItems0 struct {

	// encrypted
	// Required: true
	Encrypted *bool `json:"encrypted"`

	// modifier
	// Pattern: ^([\w=-]+)$|^((\/[^\/]+)+\/?)$
	Modifier string `json:"modifier,omitempty"`

	// operator
	// Enum: [exists equals]
	Operator string `json:"operator,omitempty"`

	// query term
	// Required: true
	// Pattern: ^([\w=-]+)$|^((\/[^\/]+)+\/?)$
	QueryTerm *string `json:"queryTerm"`

	// value
	// Max Length: 0
	Value string `json:"value,omitempty"`
}

// Validate validates this explore query panels items0 items items0
func (m *ExploreQueryPanelsItems0ItemsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEncrypted(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModifier(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOperator(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateQueryTerm(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExploreQueryPanelsItems0ItemsItems0) validateEncrypted(formats strfmt.Registry) error {

	if err := validate.Required("encrypted", "body", m.Encrypted); err != nil {
		return err
	}

	return nil
}

func (m *ExploreQueryPanelsItems0ItemsItems0) validateModifier(formats strfmt.Registry) error {

	if swag.IsZero(m.Modifier) { // not required
		return nil
	}

	if err := validate.Pattern("modifier", "body", string(m.Modifier), `^([\w=-]+)$|^((\/[^\/]+)+\/?)$`); err != nil {
		return err
	}

	return nil
}

var exploreQueryPanelsItems0ItemsItems0TypeOperatorPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["exists","equals"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		exploreQueryPanelsItems0ItemsItems0TypeOperatorPropEnum = append(exploreQueryPanelsItems0ItemsItems0TypeOperatorPropEnum, v)
	}
}

const (

	// ExploreQueryPanelsItems0ItemsItems0OperatorExists captures enum value "exists"
	ExploreQueryPanelsItems0ItemsItems0OperatorExists string = "exists"

	// ExploreQueryPanelsItems0ItemsItems0OperatorEquals captures enum value "equals"
	ExploreQueryPanelsItems0ItemsItems0OperatorEquals string = "equals"
)

// prop value enum
func (m *ExploreQueryPanelsItems0ItemsItems0) validateOperatorEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, exploreQueryPanelsItems0ItemsItems0TypeOperatorPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *ExploreQueryPanelsItems0ItemsItems0) validateOperator(formats strfmt.Registry) error {

	if swag.IsZero(m.Operator) { // not required
		return nil
	}

	// value enum
	if err := m.validateOperatorEnum("operator", "body", m.Operator); err != nil {
		return err
	}

	return nil
}

func (m *ExploreQueryPanelsItems0ItemsItems0) validateQueryTerm(formats strfmt.Registry) error {

	if err := validate.Required("queryTerm", "body", m.QueryTerm); err != nil {
		return err
	}

	if err := validate.Pattern("queryTerm", "body", string(*m.QueryTerm), `^([\w=-]+)$|^((\/[^\/]+)+\/?)$`); err != nil {
		return err
	}

	return nil
}

func (m *ExploreQueryPanelsItems0ItemsItems0) validateValue(formats strfmt.Registry) error {

	if swag.IsZero(m.Value) { // not required
		return nil
	}

	if err := validate.MaxLength("value", "body", string(m.Value), 0); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExploreQueryPanelsItems0ItemsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExploreQueryPanelsItems0ItemsItems0) UnmarshalBinary(b []byte) error {
	var res ExploreQueryPanelsItems0ItemsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}