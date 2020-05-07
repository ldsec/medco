package unlynx

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type testKeySwitchValueParameters struct {
	value int64
}

var keySwitchParameters = []testKeySwitchValueParameters{
	{
		value: 0,
	},
	{
		value: 1,
	},
	{
		value: 266,
	},
	{
		value: 9999,
	},
}

// TestUnlynx runs the unlynx tests.
func TestUnlynx() (testPassed bool) {

	for _, testParams := range keySwitchParameters {
		if !testKeySwitch(testParams) {
			log := "test failed: "
			text, err := json.Marshal(testParams)
			if err == nil {
				log += string(text)
			} else {
				log += err.Error()
			}
			logrus.Warn(log)
			return false
		}
	}

	return true
}

func testKeySwitch(testParams testKeySwitchValueParameters) (testPassed bool) {

	enc, err := EncryptWithCothorityKey(testParams.value)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	keySwitchedEnc, _, err := KeySwitchValue("test query "+time.Now().Format(time.RFC3339Nano), enc, pubKey)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	dec, err := Decrypt(keySwitchedEnc, privKey)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	if dec != testParams.value {
		return false
	}

	return true
}
