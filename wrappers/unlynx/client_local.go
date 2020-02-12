package unlynx

import (
	"github.com/ldsec/medco-connector/util/server"
	libunlynx "github.com/ldsec/unlynx/lib"
	stats "github.com/r0fls/gostats"
	"github.com/sirupsen/logrus"
)

// LocallyAggregateValues adds together several encrypted values homomorphically
// if there are no values, the value zero encrypted is returned
func LocallyAggregateValues(values []string) (agg string, err error) {

	if len(values) == 0 {
		return getEncryptedZero()
	}

	// deserialize values
	deserialized, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	// local aggregation
	aggregate :=  &deserialized[0]
	for i := 1; i < len(deserialized); i++ {
		aggregate.Add(*aggregate, deserialized[i])
	}

	agg, err = aggregate.Serialize()
	if err != nil {
		logrus.Error("unlynx error serializing: ", err)
	}
	return
}

// LocallyMultiplyScalar multiply homomorphically an encrypted value with a clear-text scalar
func LocallyMultiplyScalar(encValue string, scalar int64) (res string, err error) {

	deserialized := libunlynx.CipherText{}
	err = deserialized.Deserialize(encValue)
	if err != nil {
		logrus.Error("unlynx error deserializing cipher text: ", err)
		return
	}

	result := libunlynx.CipherText{}
	result.MulCipherTextbyScalar(deserialized, libunlynx.SuiTe.Scalar().SetInt64(scalar))

	res, err =  result.Serialize()
	if err != nil {
		logrus.Error("unlynx error serializing: ", err)
	}
	return
}

// LocallyObfuscateValue adds random noise homomorphically to an encrypted value
func LocallyObfuscateValue(encValue string, obfuscationParam int, pubKey string) (res string, err error) {

	if obfuscationParam < utilserver.MedCoObfuscationMin {
		obfuscationParam = utilserver.MedCoObfuscationMin
		logrus.Info("Obfuscation variance set to the minimum of ", obfuscationParam)
	}

	distribution := stats.Laplace(0, float64(obfuscationParam))
	noise := distribution.Random()

	// encrypt the noise
	encNoise, err := Encrypt(int64(noise), pubKey)
	if err != nil {
		return
	}

	// add together value and noise
	return LocallyAggregateValues([]string{encValue, encNoise})
}

// EncryptWithCothorityKey encrypts an integer with the public key of the cothority
func EncryptWithCothorityKey(value int64) (encrypted string, err error) {
	_, cothorityRoster := NewUnlynxClient()
	encrypted, err = libunlynx.EncryptInt(cothorityRoster.Aggregate, value).Serialize()
	if err != nil {
		logrus.Error("unlynx failed serializing encrypted value: ", err)
	}
	return
}

// Encrypt encrypts an integer with a public key
func Encrypt(value int64, pubKey string) (encrypted string, err error) {
	pubKeyDes, err := libunlynx.DeserializePoint(pubKey)
	if err != nil {
		logrus.Error("unlynx failed deserializing public key: ", err)
		return
	}

	encrypted, err = libunlynx.EncryptInt(pubKeyDes, value).Serialize()
	if err != nil {
		logrus.Error("unlynx failed serializing encrypted value: ", err)
	}
	return
}

// Decrypt decrypts an integer with a private key
func Decrypt(value string, privKey string) (decrypted int64, err error) {
	valueDes := libunlynx.CipherText{}
	err = valueDes.Deserialize(value)
	if err != nil {
		logrus.Error("unlynx error deserializing cipher text: ", err)
		return
	}

	privKeyDes, err := libunlynx.DeserializeScalar(privKey)
	if err != nil {
		logrus.Error("unlynx error deserializing scalar: ", err)
		return
	}

	decrypted = libunlynx.DecryptInt(privKeyDes, valueDes)
	return
}

// GenerateKeyPair generates a matching pair of public and private keys
func GenerateKeyPair() (pubKey string, privKey string, err error) {
	rawPrivKey, rawPubKey := libunlynx.GenKey()

	privKey, err = libunlynx.SerializeScalar(rawPrivKey)
	if err != nil {
		logrus.Error("unlynx error serializing private key: ", err)
		return
	}

	pubKey, err = libunlynx.SerializePoint(rawPubKey)
	if err != nil {
		logrus.Error("unlynx error serializing private key: ", err)
		return
	}

	return
}

// GetCothorityKey returns the aggregated cothority public key
func GetCothorityKey() (key string, err error) {
	_, cothorityRoster := NewUnlynxClient()
	key, err = libunlynx.SerializePoint(cothorityRoster.Aggregate)
	if err != nil {
		logrus.Error("unlynx error serializing public aggregate key: ", err)
	}
	return
}
