package survivalclient

import (
	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

// EncryptTimeCodes takes the map of clear text timeCodes to their integer identifier and returns
//  a map of clear text timeCodes to ElGamal encryption of the integer identifier
//  the inverse map of ElGamal encryption to clear text time point
func EncryptTimeCodes(timeCodesIntID map[string]int64) (map[string]string, map[string]string, error) {
	result := make(map[string]string, len(timeCodesIntID))
	inverseMap := make(map[string]string, len(timeCodesIntID))

	for timeValue, intID := range timeCodesIntID {
		encryptedTime, err := unlynx.EncryptWithCothorityKey(intID)
		if err != nil {
			return nil, nil, err
		}

		result[timeValue] = encryptedTime
		inverseMap[encryptedTime] = timeValue
	}
	return result, inverseMap, nil
}
