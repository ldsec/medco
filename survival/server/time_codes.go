package survivalserver

import (
	"errors"
	"time"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

type TagsToTagIDs interface {
	GetTagIDs(map[string]string) (map[string]string, error)
}

type TimeCodesMap struct {
	//thoses maps are sinks, they can only grow, and once a value is inside,
	encTimeCodesToTags   map[string]string
	encTimeCodesToTagIDs map[string]string
	tagsToEncTimeCodes   map[string]string
	tagIDsToEncTimeCodes map[string]string
	tagIDs               []string
	tagToTagIDS          TagsToTagIDs
}

func NewTimeCodesMap(queryName string, encTimeCodes []string, tagToTagIDS TagsToTagIDs) (timeCodeMap *TimeCodesMap, times map[string]time.Duration, err error) {
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
		tagToTagIDS:          tagToTagIDS,
	}
	timeCodeMap.encTimeCodesToTags, times, err = unlynx.DDTagValues(queryName+"_TIME_CONCEPT_CODES_", encTimeCodes)

	for encTimeCode, tag := range timeCodeMap.encTimeCodesToTags {
		timeCodeMap.tagsToEncTimeCodes[tag] = encTimeCode
	}

	timeCodeMap.tagIDsToEncTimeCodes, err = timeCodeMap.tagToTagIDS.GetTagIDs(timeCodeMap.tagsToEncTimeCodes)

	if err != nil {
		return
	}

	for tagID := range timeCodeMap.tagIDsToEncTimeCodes {

		timeCodeMap.tagIDs = append(timeCodeMap.tagIDs, tagID)
	}
	return

}

func (timeMap *TimeCodesMap) GetTagIDList() []string {
	return timeMap.tagIDs

}

func (timeMap *TimeCodesMap) GetTagIDMap() map[string]string {
	return timeMap.tagIDsToEncTimeCodes

}
