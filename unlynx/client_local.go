package unlynx

import (
	libunlynx "github.com/lca1/unlynx/lib"
	"github.com/sirupsen/logrus"
)

// LocallyAggregateValues adds together several encrypted values homomorphically
func LocallyAggregateValues(values []string) (agg string, err error) {

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
