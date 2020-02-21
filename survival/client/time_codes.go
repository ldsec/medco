package survivalclient

import (
	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

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
