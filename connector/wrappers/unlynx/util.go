package unlynx

import (
	"errors"
	"os"
	"strconv"

	utilserver "github.com/ldsec/medco/connector/util/server"
	servicesmedco "github.com/ldsec/medco/unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
)

// serializeCipherVector serializes a vector of cipher texts into a string-encoded slice
func serializeCipherVector(cipherVector libunlynx.CipherVector) (serializedVector []string, err error) {
	for _, cipherText := range cipherVector {

		serialized, err := cipherText.Serialize()
		if err != nil {
			logrus.Error("unlynx error serializing: ", err)
			return nil, err
		}

		serializedVector = append(serializedVector, serialized)
	}
	return
}

// deserializeCipherVector deserializes string-encoded cipher texts into a vector
func deserializeCipherVector(cipherTexts []string) (cipherVector libunlynx.CipherVector, err error) {
	for idx, cipherText := range cipherTexts {

		if len(cipherText) == 0 {
			err = errors.New("invalid value in the cipher vector index " + strconv.Itoa(idx))
			logrus.Error(err)
			return
		}

		deserialized := libunlynx.CipherText{}
		err = deserialized.Deserialize(cipherText)
		if err != nil {
			logrus.Error("unlynx error deserializing cipher text: ", err)
			return
		}

		cipherVector = append(cipherVector, deserialized)
	}
	return
}

// newUnlynxClient creates a new client to communicate with unlynx
func newUnlynxClient() (unlynxClient *servicesmedco.API, cothorityRoster *onet.Roster) {

	// initialize medco client
	groupFile, err := os.Open(utilserver.UnlynxGroupFilePath)
	if err != nil {
		logrus.Panic("unlynx error opening group file: ", err)
	}

	group, err := app.ReadGroupDescToml(groupFile)
	if err != nil || len(group.Roster.List) <= 0 {
		logrus.Panic("unlynx error parsing group file: ", err)
	}

	cothorityRoster = group.Roster
	unlynxClient = servicesmedco.NewMedCoClient(
		cothorityRoster.List[utilserver.MedCoNodeIdx],
		strconv.Itoa(utilserver.MedCoNodeIdx),
	)

	return
}

// getEncryptedZero returns an encrypted zero
func getEncryptedZero() (serializedZero string, err error) {
	_, cothorityRoster := newUnlynxClient()
	encZero := libunlynx.EncryptInt(cothorityRoster.Aggregate, 0)
	serializedZero, err = encZero.Serialize()
	if err != nil {
		logrus.Error("unlynx failed serializing zero: ", err)
	}
	return
}
