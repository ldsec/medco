// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Timing timing
//
// swagger:model timing
type Timing string

const (

	// TimingAny captures enum value "any"
	TimingAny Timing = "any"

	// TimingSamevisit captures enum value "samevisit"
	TimingSamevisit Timing = "samevisit"

	// TimingSameinstancenum captures enum value "sameinstancenum"
	TimingSameinstancenum Timing = "sameinstancenum"
)

// for schema
var timingEnum []interface{}

func init() {
	var res []Timing
	if err := json.Unmarshal([]byte(`["any","samevisit","sameinstancenum"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		timingEnum = append(timingEnum, v)
	}
}

func (m Timing) validateTimingEnum(path, location string, value Timing) error {
	if err := validate.EnumCase(path, location, value, timingEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this timing
func (m Timing) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateTimingEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
