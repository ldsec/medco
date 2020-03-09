package survivalserver

import (
	"errors"
	"fmt"
	"strings"
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

func (timeCodesMap *TimeCodesMap) getTagIDs() (err error) {
	paths := buildParameters(timeCodesMap.tagsToEncTimeCodes)
	psqlQuery := tagQuery + ` WHERE concept_path IN (` + paths + `)`
	rows, err := directAccessDB.Query(psqlQuery)
	err = NiceError(err)
	if err != nil {
		return
	}
	timeCodesMap.tagIDsToEncTimeCodes = make(map[string]string, len(timeCodesMap.tagsToEncTimeCodes))
	numberOfRows := 0

	for rows.Next() {
		var path, conceptCode string
		err = rows.Scan(&path, &conceptCode)

		if err != nil {
			return
		}
		tag := strings.Replace(path, "\\medco\\tagged\\concept\\", "", 1)
		tag = strings.Replace(tag, "\\", "", 1)
		encTimeCode, ok := timeCodesMap.tagsToEncTimeCodes[tag]
		if !ok {
			err = fmt.Errorf("tag  %s  not found, node index %d", tag, utilserver.MedCoNodeIdx)
			return
		}

		timeCodesMap.tagIDsToEncTimeCodes[conceptCode] = encTimeCode
		numberOfRows++

	}
	if numberOfRows == 0 {
		err = fmt.Errorf("From node %d, Unable to find any of the tag in the data base %s", utilserver.MedCoNodeIdx, dbName)
	}

	err = rows.Close()

	return
}

func buildParameters(tags map[string]string) string {
	paths := make([]string, len(tags))
	pos := 0
	for tag := range tags {
		paths[pos] = `'\medco\tagged\concept\` + tag + `\'`
		pos++
	}

	return strings.Join(paths, ",")
}

//TimeCodesMap holds the different mappings within encrypted time codes, time code tags and time code tag identifier
type TimeCodesMap struct {
	//thoses maps are sinks, they can only grow, and once a value is inside,
	encTimeCodesToTags   map[string]string
	encTimeCodesToTagIDs map[string]string
	tagsToEncTimeCodes   map[string]string
	tagIDsToEncTimeCodes map[string]string
	tagIDs               []string
}

//NewTimeCodesMap time codes map constructor, it implies requests to the database and the unlynx module
func NewTimeCodesMap(queryName string, encTimeCodes []string) (timeCodeMap *TimeCodesMap, times map[string]time.Duration, err error) {
	length := len(encTimeCodes)
	if length == 0 {
		err = errors.New("Empty list of time codes")
		return
	}
	timeCodeMap = &TimeCodesMap{
		encTimeCodesToTags:   make(map[string]string, length),
		encTimeCodesToTagIDs: make(map[string]string, length),
		tagsToEncTimeCodes:   make(map[string]string, length),
		tagIDsToEncTimeCodes: make(map[string]string, length),
		tagIDs:               make([]string, 0, length),
	}
	timeCodeMap.encTimeCodesToTags, times, err = unlynx.DDTagValues(queryName+"_TIME_CONCEPT_CODES_", encTimeCodes)

	for encTimeCode, tag := range timeCodeMap.encTimeCodesToTags {
		timeCodeMap.tagsToEncTimeCodes[tag] = encTimeCode
	}

	err = timeCodeMap.getTagIDs()

	if err != nil {
		return
	}

	for tagID := range timeCodeMap.tagIDsToEncTimeCodes {

		timeCodeMap.tagIDs = append(timeCodeMap.tagIDs, tagID)
	}
	return

}

// GetTagIDList returns the tag ID list
func (timeCodesMap *TimeCodesMap) GetTagIDList() []string {
	return timeCodesMap.tagIDs

}

//GetTagIDMap returns the mapping from tag identifiers to time codes encrypted under collective authority key
func (timeCodesMap *TimeCodesMap) GetTagIDMap() map[string]string {
	return timeCodesMap.tagIDsToEncTimeCodes

}
