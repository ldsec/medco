package explorestatisticsclient

import "github.com/ldsec/medco/connector/restapi/models"

// NodeClearResults holds all needed information to create an histogram, in clear text.
// EncryptedResults is the mirror structure of NodeClearResults, before the decryption.
// NodeClearResults implements sort.Interface interface.
type NodeClearResults struct {
	Intervals []*models.ClearInterval
	Unit      string
}

// // Len implements Len method for sort.Interface interface
// func (res ClearResults) Len() int {
// 	return len(res)
// }

// // Less implements Less method for sort.Interface interface
// func (res ClearResults) Less(i, j int) bool {
// 	return res[i].GroupID < res[j].GroupID
// }

// //Swap implements Swap method for sort.Interface interface
// func (res ClearResults) Swap(i, j int) {
// 	res[i], res[j] = res[j], res[i]
// }
