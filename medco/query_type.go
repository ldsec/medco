package medco

import "github.com/lca1/medco-connector/swagger/models"

// I2b2MedCoQueryType encodes the requested results from the query
type I2b2MedCoQueryType struct {

	// PatientList encodes if the (encrypted) patient list is retrieved, otherwise only a count
	PatientList bool

	// CountPerSite encodes if the retrieved count is per site, or aggregated over all nodes
	CountPerSite bool

	// Shuffled encodes if the count per site if shuffled across all nodes for unlinkability
	Shuffled bool

	// Obfuscated encodes if the count (per site or aggregated) should be obfuscated
	Obfuscated bool

	// todo: include differential privacy requested
}


// resolveQueryType resolves the query type based on the user's authorizations and its request
func resolveQueryType(query *models.QueryI2b2Medco, userPrincipal *models.User) (queryType I2b2MedCoQueryType) {
	// todo: authorizations param model to create
	// todo: other authorizations

	queryType = I2b2MedCoQueryType{
		PatientList: false,
		CountPerSite: true, // todo: from principal
		Shuffled: true, // todo: from principal
		Obfuscated: false, // todo: with DP implementation
	}

	for _, selectedResult := range query.Select {
		if selectedResult == "patient_list" {
			queryType.PatientList = true
		}
	}

	return
}
