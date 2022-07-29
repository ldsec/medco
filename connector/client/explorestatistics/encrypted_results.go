package explorestatisticsclient

import (
	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
)

// EncryptedResults holds all needed information to perform survival analysis, encrypted with the client public key.
// ClearResults is the mirror structure of EncryptedResults, once decrpyted.
type EncryptedResults struct {
	Intervals []*models.IntervalBucket
	Unit      string
	Timers    medcomodels.Timers
}

// Decrypt deciphers initial counts, numbers of censoring events and events of interest with the provided key
func (nodeResults EncryptedResults) Decrypt(privateKey string) (*NodeClearResults, error) {
	clrIntervals := make([]*models.ClearInterval, len(nodeResults.Intervals))
	for idxInterval, encInterval := range nodeResults.Intervals {
		count, err := unlynx.Decrypt(*encInterval.EncCount, privateKey)
		if err != nil {
			return nil, err
		}
		clrInterval := models.ClearInterval{
			Count:       &count,
			HigherBound: encInterval.HigherBound,
			LowerBound:  encInterval.LowerBound,
		}
		clrIntervals[idxInterval] = &clrInterval
	}

	clrResults := new(NodeClearResults)
	clrResults.Intervals = clrIntervals
	clrResults.Unit = nodeResults.Unit
	return clrResults, nil
}
