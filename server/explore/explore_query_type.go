package medcoserver

import (
	"errors"

	"github.com/ldsec/medco-connector/restapi/models"
)

// ExploreQueryType encodes the requested results from the query
type ExploreQueryType struct {

	// PatientList encodes if the (encrypted) patient list is retrieved, otherwise only a count
	PatientList bool

	// CountType encodes the type of count, either global or per site
	CountType CountType

	// Shuffled encodes if the count per site if shuffled across all nodes for unlinkability
	Shuffled bool

	// Obfuscated encodes if the count (per site or aggregated) should be obfuscated
	Obfuscated bool

	// todo: include differential privacy requested

	// PatientSet  encodes if the patient set are needed further
	PatientSet bool
}

// CountType encodes the type of count, either global or per site
type CountType int

const (
	// PerSite represents the count broken down per site
	PerSite = iota

	// Global represents the count globally aggregated over all sites
	Global
)

// NewExploreQueryType creates a query type based on the requested query type
func NewExploreQueryType(requestedQueryType models.ExploreQueryType) (queryType ExploreQueryType, err error) {

	switch requestedQueryType {
	case models.ExploreQueryTypePatientList:
		queryType = ExploreQueryType{
			PatientList: true,
			CountType:   PerSite,
			Shuffled:    false,
			Obfuscated:  false,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountPerSite:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   PerSite,
			Shuffled:    false,
			Obfuscated:  false,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountPerSiteObfuscated:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   PerSite,
			Shuffled:    false,
			Obfuscated:  true,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountPerSiteShuffled:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   PerSite,
			Shuffled:    true,
			Obfuscated:  false,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountPerSiteShuffledObfuscated:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   PerSite,
			Shuffled:    true,
			Obfuscated:  true,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountGlobal:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   Global,
			Shuffled:    false,
			Obfuscated:  false,
			PatientSet:  false,
		}

	case models.ExploreQueryTypeCountGlobalObfuscated:
		queryType = ExploreQueryType{
			PatientList: false,
			CountType:   Global,
			Shuffled:    false,
			Obfuscated:  true,
			PatientSet:  false,
		}

	default:
		err = errors.New("unrecognized query type: " + string(requestedQueryType))
	}

	return
}
