// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// ExploreQueryType explore query type
// swagger:model exploreQueryType
type ExploreQueryType string

const (

	// ExploreQueryTypePatientList captures enum value "patient_list"
	ExploreQueryTypePatientList ExploreQueryType = "patient_list"

	// ExploreQueryTypeCountPerSite captures enum value "count_per_site"
	ExploreQueryTypeCountPerSite ExploreQueryType = "count_per_site"

	// ExploreQueryTypeCountPerSiteObfuscated captures enum value "count_per_site_obfuscated"
	ExploreQueryTypeCountPerSiteObfuscated ExploreQueryType = "count_per_site_obfuscated"

	// ExploreQueryTypeCountPerSiteShuffled captures enum value "count_per_site_shuffled"
	ExploreQueryTypeCountPerSiteShuffled ExploreQueryType = "count_per_site_shuffled"

	// ExploreQueryTypeCountPerSiteShuffledObfuscated captures enum value "count_per_site_shuffled_obfuscated"
	ExploreQueryTypeCountPerSiteShuffledObfuscated ExploreQueryType = "count_per_site_shuffled_obfuscated"

	// ExploreQueryTypeCountGlobal captures enum value "count_global"
	ExploreQueryTypeCountGlobal ExploreQueryType = "count_global"

	// ExploreQueryTypeCountGlobalObfuscated captures enum value "count_global_obfuscated"
	ExploreQueryTypeCountGlobalObfuscated ExploreQueryType = "count_global_obfuscated"

	// ExploreQueryTypeSurvivalAnalysis captures enum value "survival_analysis"
	ExploreQueryTypeSurvivalAnalysis ExploreQueryType = "survival_analysis"
)

// for schema
var exploreQueryTypeEnum []interface{}

func init() {
	var res []ExploreQueryType
	if err := json.Unmarshal([]byte(`["patient_list","count_per_site","count_per_site_obfuscated","count_per_site_shuffled","count_per_site_shuffled_obfuscated","count_global","count_global_obfuscated","survival_analysis"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		exploreQueryTypeEnum = append(exploreQueryTypeEnum, v)
	}
}

func (m ExploreQueryType) validateExploreQueryTypeEnum(path, location string, value ExploreQueryType) error {
	if err := validate.Enum(path, location, value, exploreQueryTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this explore query type
func (m ExploreQueryType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateExploreQueryTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
