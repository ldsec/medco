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
	psqlQuery := tagQuery + ` WHERE time_path IN (` + paths + `)`
	rows, err := directAccessDB.Query(psqlQuery)
	err = NiceError(err)
	if err != nil {
		return
	}
	timeCodesMap.tagIDsToEncTimeCodes = make(map[TagID]EncryptedEncID, len(timeCodesMap.tagsToEncTimeCodes))
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

		timeCodesMap.tagIDsToEncTimeCodes[TagID(conceptCode)] = encTimeCode
		numberOfRows++

	}
	if numberOfRows == 0 {
		err = fmt.Errorf("From node %d, Unable to find any of the tag in the data base ", utilserver.MedCoNodeIdx)
		return
	}

	err = rows.Close()

	return
}

func buildParameters(tags map[string]EncryptedEncID) string {
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
	encTimeCodesToTags   map[EncryptedEncID]string
	encTimeCodesToTagIDs map[EncryptedEncID]TagID
	tagsToEncTimeCodes   map[string]EncryptedEncID
	tagIDsToEncTimeCodes map[TagID]EncryptedEncID
	tagIDs               []TagID
}

//NewTimeCodesMap time codes map constructor, it implies requests to the database and the unlynx module
func NewTimeCodesMap(queryName string, encTimeCodes []EncryptedEncID) (timeCodeMap *TimeCodesMap, times map[string]time.Duration, err error) {
	length := len(encTimeCodes)
	if length == 0 {
		err = errors.New("Empty list of time codes")
		return
	}
	timeCodeMap = &TimeCodesMap{
		encTimeCodesToTags:   make(map[EncryptedEncID]string, length),
		encTimeCodesToTagIDs: make(map[EncryptedEncID]TagID, length),
		tagsToEncTimeCodes:   make(map[string]EncryptedEncID, length),
		tagIDsToEncTimeCodes: make(map[TagID]EncryptedEncID, length),
		tagIDs:               make([]TagID, 0, length),
	}
	encTimeCodesString := make([]string, length)
	for idx, timeCode := range encTimeCodes {
		encTimeCodesString[idx] = string(timeCode)
	}
	encTimeCodesToTagsString, times, err := unlynx.DDTagValues(queryName+"_TIME_CONCEPT_CODES_", encTimeCodesString)
	for timeCode, tag := range encTimeCodesToTagsString {
		timeCodeMap.encTimeCodesToTags[EncryptedEncID(timeCode)] = tag
	}
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

//NewTimeCodesMapWithCallback has the same purpose of NewTimeCodesMap, but return an error chan as the error can occur after the call of this function
func NewTimeCodesMapWithCallback(queryName string, encTimeCodes []EncryptedEncID, callBack func(*TimeCodesMap, chan error, chan map[string]time.Duration)) (<-chan *TimeCodesMap, <-chan map[string]time.Duration, <-chan error) {
	timeCodeMapChan := make(chan *TimeCodesMap, 1)
	timesChan := make(chan map[string]time.Duration, 1)
	errChan := make(chan error, 1)
	interMediateTimesChan := make(chan map[string]time.Duration, 1)

	errChan = make(chan error, 1) //if not 1, the fact of extrating these channels in a undefined order would block at pushing (the first possibility is immediate with errChan)
	length := len(encTimeCodes)
	if length == 0 {
		errChan <- errors.New("Empty list of time codes")
		return timeCodeMapChan, timesChan, errChan
	}

	encTimeCodesStrings := make([]string, length)

	for idx, timeCode := range encTimeCodes {
		encTimeCodesStrings[idx] = string(timeCode)
	}

	go func() {
		timer := make(map[string]time.Duration)
		timeCodeMap := &TimeCodesMap{
			encTimeCodesToTags:   make(map[EncryptedEncID]string, length),
			encTimeCodesToTagIDs: make(map[EncryptedEncID]TagID, length),
			tagsToEncTimeCodes:   make(map[string]EncryptedEncID, length),
			tagIDsToEncTimeCodes: make(map[TagID]EncryptedEncID, length),
			tagIDs:               make([]TagID, 0, length),
		}
		encTimeCodesToTagsStrings, times, err := unlynx.DDTagValues(queryName+"_TIME_CONCEPT_CODES_", encTimeCodesStrings)
		for key, value := range times {
			timer[key] = value
		}
		for timeCode, tag := range encTimeCodesToTagsStrings {
			timeCodeMap.encTimeCodesToTags[EncryptedEncID(timeCode)] = tag
		}
		for encTimeCode, tag := range timeCodeMap.encTimeCodesToTags {
			timeCodeMap.tagsToEncTimeCodes[tag] = encTimeCode
		}

		err = timeCodeMap.getTagIDs()

		if err != nil {
			errChan <- err
		}

		for tagID := range timeCodeMap.tagIDsToEncTimeCodes {

			timeCodeMap.tagIDs = append(timeCodeMap.tagIDs, tagID)
		}

		callBack(timeCodeMap, errChan, interMediateTimesChan)
		select {
		case intermediateTimer := <-interMediateTimesChan:
			for key, val := range intermediateTimer {
				timer[key] = val
			}
		default:
		}

		timeCodeMapChan <- timeCodeMap
		timesChan <- timer
		return
	}()

	return timeCodeMapChan, timesChan, errChan
}

// GetTagIDList returns the tag ID list
func (timeCodesMap *TimeCodesMap) GetTagIDList() []TagID {
	return timeCodesMap.tagIDs

}

//GetTagIDMap returns the mapping from tag identifiers to time codes encrypted under collective authority key
func (timeCodesMap *TimeCodesMap) GetTagIDMap() map[TagID]EncryptedEncID {
	return timeCodesMap.tagIDsToEncTimeCodes

}
