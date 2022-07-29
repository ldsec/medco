package explorestatisticsclient

// import (
// 	"testing"

// 	medcomodels "github.com/ldsec/medco/connector/models"

// 	"github.com/stretchr/testify/assert"
// )

// const testPubKey = "mNd1d91GTy4wI0ZNgN7pTPo89GG7d_oOrGCPRUNG4gE="
// const testPrivateKey = "NMfXoUyb1yRSY3goHboUK9HEbMSvuzEmQTRHO_eZSAk="

// func TestDecryptGroup(t *testing.T) {
// 	decrypted, err := testEncryptedResults.Decrypt(testPrivateKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, testDecryptedResults, decrypted)
// 	corrupted := testEncryptedResults
// 	corrupted[0].EncryptedInitialCount = "This is not a cipher text"
// 	_, err = testEncryptedResults.Decrypt(testPrivateKey)
// 	assert.Error(t, err)
// 	corrupted1 := testEncryptedResults
// 	corrupted1[0].TimePoints[0].EncryptedEvents.CensoringEvents = "This is not a cipher text"
// 	_, err = testEncryptedResults.Decrypt(testPrivateKey)
// 	assert.Error(t, err)

// }

// var testEncryptedResults = EncryptedResults{
// 	{GroupID: "testGroup1",
// 		EncryptedInitialCount: "YKu4hdlub0k7VKrHmVEMwDTnNLEcuHypqrhTfvXK9ABDGR4f1jw7vHhO3jNViQI4I-W8tGu8G2FpjamnEyN1OA==",
// 		TimePoints: []struct {
// 			Time            int
// 			EncryptedEvents struct {
// 				EventsOfInterest string
// 				CensoringEvents  string
// 			}
// 		}{
// 			{
// 				Time: 1,
// 				EncryptedEvents: struct {
// 					EventsOfInterest string
// 					CensoringEvents  string
// 				}{
// 					EventsOfInterest: "EjgEe_xmWLpBactWR_Dyl1EIkMtLpc0P8Zn2aTT3VQ0GXTdwgg8jjNW6J0EgGxwWKrXoyDi31w2TXLs_W9GxCQ==",
// 					CensoringEvents:  "9_lAY7H_u0AUiWeFW8B3J_5EqyWBfyEGYcxTfq152jfunwDvbCcVXtqGCUkeb2wgbHxRoCEe93RSVcm_Z8wstg==",
// 				},
// 			}},
// 	},
// }
// var testDecryptedResults = ClearResults{
// 	{GroupID: "testGroup1",
// 		InitialCount: 50,
// 		TimePoints: medcomodels.TimePoints{
// 			{
// 				Time: 1,
// 				Events: struct {
// 					EventsOfInterest int64
// 					CensoringEvents  int64
// 				}{
// 					EventsOfInterest: 10,
// 					CensoringEvents:  0,
// 				},
// 			}},
// 	},
// }
