package survivalclient

import (
	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"
	"github.com/ldsec/medco/connector/util"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
)

// EncryptedResults holds all needed information to perform survival analysis, encrypted with the client public key.
// ClearResults is the mirror structure of EncryptedResults, once decrpyted.
type EncryptedResults []struct {
	GroupID               string
	EncryptedInitialCount string
	TimePoints            []struct {
		Time            int
		EncryptedEvents struct {
			EventsOfInterest string
			CensoringEvents  string
		}
	}
}

// Decrypt deciphers initial counts, numbers of censoring events and events of interest with the provided key
func (ciphers EncryptedResults) Decrypt(privateKey string) (ClearResults, error) {
	res := make(ClearResults, len(ciphers))
	var err error
	for i, cipher := range ciphers {
		res[i].GroupID = cipher.GroupID
		res[i].InitialCount, err = unlynx.Decrypt(cipher.EncryptedInitialCount, privateKey)
		if err != nil {
			return nil, err
		}
		res[i].TimePoints = make(util.TimePoints, len(cipher.TimePoints))
		for j, events := range cipher.TimePoints {
			res[i].TimePoints[j].Time = events.Time
			res[i].TimePoints[j].Events.EventsOfInterest, err = unlynx.Decrypt(events.EncryptedEvents.EventsOfInterest, privateKey)
			if err != nil {
				return nil, err
			}
			res[i].TimePoints[j].Events.CensoringEvents, err = unlynx.Decrypt(events.EncryptedEvents.CensoringEvents, privateKey)
			if err != nil {
				return nil, err
			}

		}
	}
	return res, nil
}

// encryptedResultsFromAPIResponse converts a survival analysis API response into a EncryptedResults instance
func encryptedResultsFromAPIResponse(bodyResults []*survival_analysis.SurvivalAnalysisOKBodyResultsItems0) EncryptedResults {
	res := make(EncryptedResults, len(bodyResults))
	for i, group := range bodyResults {
		res[i].GroupID = group.GroupID
		res[i].EncryptedInitialCount = group.InitialCount
		res[i].TimePoints = make([]struct {
			Time            int
			EncryptedEvents struct {
				EventsOfInterest string
				CensoringEvents  string
			}
		}, len(group.GroupResults))
		for j, timePoint := range group.GroupResults {
			res[i].TimePoints[j].Time = int(timePoint.Timepoint)
			res[i].TimePoints[j].EncryptedEvents.EventsOfInterest = timePoint.Events.Eventofinterest
			res[i].TimePoints[j].EncryptedEvents.CensoringEvents = timePoint.Events.Censoringevent
		}
	}
	return res
}
