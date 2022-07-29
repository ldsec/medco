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

// RestAPIAuthorization rest Api authorization
//
// swagger:model restApiAuthorization
type RestAPIAuthorization string

const (

	// RestAPIAuthorizationMedcoNetwork captures enum value "medco-network"
	RestAPIAuthorizationMedcoNetwork RestAPIAuthorization = "medco-network"

	// RestAPIAuthorizationMedcoExplore captures enum value "medco-explore"
	RestAPIAuthorizationMedcoExplore RestAPIAuthorization = "medco-explore"

	// RestAPIAuthorizationMedcoGenomicAnnotations captures enum value "medco-genomic-annotations"
	RestAPIAuthorizationMedcoGenomicAnnotations RestAPIAuthorization = "medco-genomic-annotations"

	// RestAPIAuthorizationMedcoSurvivalAnalysis captures enum value "medco-survival-analysis"
	RestAPIAuthorizationMedcoSurvivalAnalysis RestAPIAuthorization = "medco-survival-analysis"

	// RestAPIAuthorizationMedcoExploreStatistics captures enum value "medco-explore-statistics"
	RestAPIAuthorizationMedcoExploreStatistics RestAPIAuthorization = "medco-explore-statistics"
)

// for schema
var restApiAuthorizationEnum []interface{}

func init() {
	var res []RestAPIAuthorization
	if err := json.Unmarshal([]byte(`["medco-network","medco-explore","medco-genomic-annotations","medco-survival-analysis","medco-explore-statistics"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		restApiAuthorizationEnum = append(restApiAuthorizationEnum, v)
	}
}

func (m RestAPIAuthorization) validateRestAPIAuthorizationEnum(path, location string, value RestAPIAuthorization) error {
	if err := validate.EnumCase(path, location, value, restApiAuthorizationEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this rest Api authorization
func (m RestAPIAuthorization) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateRestAPIAuthorizationEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
