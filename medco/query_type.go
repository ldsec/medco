package medco

import (
	"errors"
	"github.com/lca1/medco-connector/models"
)

// I2b2MedCoQueryType encodes the requested results from the query
type I2b2MedCoQueryType struct {

	// PatientList encodes if the (encrypted) patient list is retrieved, otherwise only a count
	PatientList bool

	// CountType encodes the type of count, either global or per site
	CountType CountType

	// Shuffled encodes if the count per site if shuffled across all nodes for unlinkability
	Shuffled bool

	// Obfuscated encodes if the count (per site or aggregated) should be obfuscated
	Obfuscated bool

	// todo: include differential privacy requested
}

// CountType encodes the type of count, either global or per site
type CountType int
const (
	// PerSite represents the count broken down per site
	PerSite = iota

	// Global represents the count globally aggregated over all sites
	Global
)


// NewI2b2MedCoQueryType creates a query type based on the requested query type
func NewI2b2MedCoQueryType(requestedQueryType models.QueryType) (queryType I2b2MedCoQueryType, err error) {

	switch requestedQueryType {
	case models.QueryTypePatientList:
		queryType = I2b2MedCoQueryType{
			PatientList: true,
			CountType: PerSite,
			Shuffled: false,
			Obfuscated: false,
		}

	case models.QueryTypeCountPerSite:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: PerSite,
			Shuffled: false,
			Obfuscated: false,
		}

	case models.QueryTypeCountPerSiteObfuscated:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: PerSite,
			Shuffled: false,
			Obfuscated: true,
		}

	case models.QueryTypeCountPerSiteShuffled:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: PerSite,
			Shuffled: true,
			Obfuscated: false,
		}

	case models.QueryTypeCountPerSiteShuffledObfuscated:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: PerSite,
			Shuffled: true,
			Obfuscated: true,
		}

	case models.QueryTypeCountGlobal:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: Global,
			Shuffled: false,
			Obfuscated: false,
		}

	case models.QueryTypeCountGlobalObfuscated:
		queryType = I2b2MedCoQueryType{
			PatientList: false,
			CountType: Global,
			Shuffled: false,
			Obfuscated: true,
		}

	default:
		err = errors.New("unrecognized query type: " + string(requestedQueryType))
	}

	return
}
