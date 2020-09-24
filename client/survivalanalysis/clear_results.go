package survivalclient

import (
	utilcommon "github.com/ldsec/medco-connector/util/common"
)

// ClearResults holds all needed information to perform survival analysis, in clear text.
// EncryptedResults is the mirror structure of ClearResults, before the encryption.
// ClearResults implements sort.Interface interface.
type ClearResults []struct {
	GroupID      string
	InitialCount int64
	TimePoints   utilcommon.TimePoints
}

// Len implements Len method for sort.Interface interface
func (res ClearResults) Len() int {
	return len(res)
}

// Less implements Less method for sort.Interface interface
func (res ClearResults) Less(i, j int) bool {
	return res[i].GroupID < res[j].GroupID
}

//Swap implements Swap method for sort.Interface interface
func (res ClearResults) Swap(i, j int) {
	res[i], res[j] = res[j], res[i]
}
