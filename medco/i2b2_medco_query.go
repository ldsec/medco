package medco

import (
	"errors"
	"github.com/lca1/medco-connector/i2b2"
	"github.com/lca1/medco-connector/swagger/models"
	"github.com/lca1/medco-connector/unlynx"
	"github.com/sirupsen/logrus"
)

// I2b2MedCoQuery executes an i2b2-medco query
func I2b2MedCoQuery(queryName string, query *models.QueryI2b2Medco) (result *models.QueryResultElement, err error) {
	// todo: timers
	// todo: log query (with associated status)
	// todo: query types (agregated per site, obfuscated per site, aggregated total)
	// todo: query type: patient list
	// todo: put user + query type + unique ID in query name

	// tag query terms
	taggedQueryTerms, err := unlynx.GetQueryTermsDDT(queryName, extractEncryptedQueryTerms(query.Panels))
	if err != nil {
		return
	}
	logrus.Info(queryName, ": tagged ", len(taggedQueryTerms), " elements with unlynx")

	// i2b2 PSM query with tagged items
	panelsItemKeys, panelsIsNot, err := prepareI2b2PsmQuery(query.Panels, taggedQueryTerms)
	if err != nil {
		return
	}

	patientCount, patientSetID, err := i2b2.ExecutePsmQuery(queryName, panelsItemKeys, panelsIsNot)
	if err != nil {
		return
	}
	logrus.Info(queryName, ": got ", patientCount, " in patient set ", patientSetID, " with i2b2")

	// i2b2 PDO query to get the dummy flags
	patientIDs, patientDummyFlags, err := i2b2.GetPatientSet(patientSetID)
	if err != nil {
		return
	}
	logrus.Info(queryName, ": got ", len(patientIDs), " patient IDs and ", len(patientDummyFlags), " dummy flags with i2b2")

	// aggregate and key-switch the result
	encCount, err := unlynx.AggregateAndKeySwitchDummyFlags(queryName, patientDummyFlags, query.UserPublicKey)

	result = &models.QueryResultElement{
		EncryptedCount: encCount,
		EncryptionKey: query.UserPublicKey,
		PatientList: nil, // todo: when implementing query type
	}
	return
}

func extractEncryptedQueryTerms(panels []*models.QueryI2b2MedcoPanelsItems0) (encQueryTerms []string) {
	for _, panel := range panels {
		for _, item := range panel.Items {
			if *item.Encrypted {
				encQueryTerms = append(encQueryTerms, item.QueryTerm)
			}
		}
	}
	return
}

func prepareI2b2PsmQuery(panels []*models.QueryI2b2MedcoPanelsItems0, taggedQueryTerms map[string]string) (panelsItemKeys [][]string, panelsIsNot []bool, err error) {
	for panelIdx, panel := range panels {
		panelsIsNot = append(panelsIsNot, *panel.Not)

		panelsItemKeys = append(panelsItemKeys, []string{})
		for _, item := range panel.Items {
			var itemKey string
			if *item.Encrypted {

				if tag, ok := taggedQueryTerms[item.QueryTerm]; ok {
					itemKey = `\\SENSITIVE_TAGGED\medco\tagged\` + tag + `\`
				} else {
					err = errors.New("query error: encrypted term does not have corresponding tag")
					logrus.Error(err)
					return
				}

			} else {
				itemKey =  item.QueryTerm
			}
			panelsItemKeys[panelIdx] = append(panelsItemKeys[panelIdx], itemKey)
		}
	}
	return
}