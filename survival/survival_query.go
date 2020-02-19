package survival

import (
	"errors"
	"strconv"
	"time"

	"github.com/ldsec/medco-connector/restapi/models"
	medcoserver "github.com/ldsec/medco-connector/server"
	survivalserver "github.com/ldsec/medco-connector/survival/server"
	"github.com/ldsec/medco-connector/survival/server/directaccess"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

//TODO remove this magic
const queryname = `I'm the survival query !`

type Query struct {
	ID            string
	UserPublicKey string
	Panels        []*models.ExploreQueryPanelsItems0

	encQueryTerm *[]string
	timeCodes    *[]string
	patientList  *[]string //nil means not yet computed, null string means computed but returned empty list
	//TODO also hide that
	Result *struct {
		EncCount string
		//EncPatientList []string
		Timers    map[string]time.Duration
		EncEvents map[string][2]string
	}

	spin     *survivalserver.Spin
	analogon *medcoserver.ExploreQuery
}

func NewQuery(qID, pubKey string, panels []*models.ExploreQueryPanelsItems0) *Query {
	return &Query{ID: qID, UserPublicKey: pubKey, Panels: panels, spin: survivalserver.NewSpin()}
}

func (q *Query) LoadTimeCodes(granularity string) error {
	q.spin.Lock()
	if q.timeCodes == nil {
		q.spin.Unlock()
		return errors.New("time codes already loaded")
	}
	q.spin.Unlock()
	timeCodes, err := getTimeCodes(granularity)
	if err != nil {

		return err
	}
	q.spin.Lock()
	q.timeCodes = &timeCodes
	q.spin.Unlock()

	return nil
}
func (query *Query) LoadPatients() (err error) {
	query.spin.Lock()
	if query.analogon == nil {
		query.spin.Unlock()
		err = errors.New("The survival query terms not converted into explore query for loading patients from panels")

		return
	}
	q := *query.analogon
	query.spin.Unlock()
	taggedQueryTerms, _, err := unlynx.DDTagValues(q.ID, q.GetEncQueryTerms())
	if err != nil {

		return
	}

	panelsItemKeys, panelsIsNot, err := q.GetI2b2PsmQueryTerms(taggedQueryTerms)
	if err != nil {

		return
	}

	_, patientSetID, err := i2b2.ExecutePsmQuery(q.ID, panelsItemKeys, panelsIsNot)
	if err != nil {

		return
	}

	patientIDs, patientDummyFlags, err := i2b2.GetPatientSet(patientSetID)
	if err != nil {

		return
	}

	aggPatientFlags, err := unlynx.LocallyAggregateValues(patientDummyFlags)
	if err != nil {

		return
	}
	//var encCount string
	//var ksCountTimers map[string]time.Duration
	logrus.Info(q.ID, ": global aggregate requested")
	_, _, err = unlynx.AggregateAndKeySwitchValue(q.ID, aggPatientFlags, q.Query.UserPublicKey)
	query.spin.Lock()
	query.patientList = &patientIDs
	query.spin.Unlock()
	return
}

func (q *Query) GetPatients() (pList []string, err error) {
	if q.patientList == nil {
		err = errors.New("patient list not loaded")
		return
	}
	pList = *(q.patientList)
	return
}

func (q *Query) GetTimeCodes() (tcd []string, err error) {
	if q.timeCodes == nil {
		err = errors.New("time codes not loaded")
		return
	}
	tcd = *(q.timeCodes)
	return
}

//like ExploreQuery

func (q *Query) convert() (err error) {
	newModelQuery := &models.ExploreQuery{Panels: q.Panels, UserPublicKey: q.UserPublicKey}
	q.analogon = new(medcoserver.ExploreQuery)
	*q.analogon, err = medcoserver.NewExploreQuery(q.ID, newModelQuery)
	return

}

func getTimeCodes(granularity string) (timeCode []string, err error) {
	results, err := i2b2.GetOntologyChildren(survivalserver.TimeConceptRootPath + "/" + granularity + "/")
	if err != nil {
		return
	}
	for _, result := range results {
		timeCode = append(timeCode, strconv.FormatInt(*(result.MedcoEncryption.ID), 10))
	}
	return

}

func (q *Query) GetID() string {
	return q.ID
}

func (q *Query) GetUserPublicKey() string {
	return q.UserPublicKey
}

func (q *Query) SetResultMap(resultMap map[string][2]string) error {
	q.spin.Lock()
	defer q.spin.Unlock()
	if q.Result == nil {
		q.Result = new(struct {
			EncCount  string
			Timers    map[string]time.Duration
			EncEvents map[string][2]string
		})
	}
	q.Result.EncEvents = resultMap
	return nil

}

func (q *Query) Execute() error {
	return directaccess.QueryTimePoints(q, 1)
}

/*
func Init() {
	survivalserver.ExecCallback = directaccess.QueryTimePoints
}
*/

//like similar to explore_query_logic, but is not necessery to share the patient list,only the count

//this one must be run only by one node

/*
func getConceptInteger{

}


func getPatient(){
	queryType:=
	queryString:=
}

//get patientList from pannels

func GetPatientList(token, username, password, queryType, queryString string, disableTLSCheck bool) (patientList []string, err error) {
	// get token
	var accessToken string
	if len(token) > 0 {
		accessToken = token
	} else {
		logrus.Debug("No token provided, requesting token for user ", username, ", disable TLS check: ", disableTLSCheck)
		accessToken, err = utilclient.RetrieveAccessToken(username, password, disableTLSCheck)
		if err != nil {
			return
		}
	}

	// parse query type
	queryTypeParsed := models.ExploreQueryType(queryType)
	err = queryTypeParsed.Validate(nil)
	if err != nil {
		logrus.Error("invalid query type")
		return
	}

	// parse query string
	panelsItemKeys, panelsIsNot, err := medcoclient.ParseQueryString(queryString)
	if err != nil {
		return
	}

	// encrypt the item keys
	encPanelsItemKeys := make([][]string, 0)
	for _, panel := range panelsItemKeys {
		encItemKeys := make([]string, 0)
		for _, itemKey := range panel {
			var encrypted string
			encrypted, err = unlynx.EncryptWithCothorityKey(itemKey)
			if err != nil {
				return
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanelsItemKeys = append(encPanelsItemKeys, encItemKeys)
	}

	// execute query
	clientQuery, err := medcoclient.NewExploreQuery(accessToken, queryTypeParsed, encPanelsItemKeys, panelsIsNot, disableTLSCheck)
	if err != nil {
		return
	}
}
*/

//get conceptList for time ID
