package utilclient

import "github.com/ldsec/medco/connector/wrappers/unlynx"

// EncryptPanelItemKeys encrypts the code the item keys with the cothority key
func EncryptPanelItemKeys(panelsItemKeys [][]int64) ([][]string, error) {
	encPanelsItemKeys := make([][]string, 0)
	for _, panel := range panelsItemKeys {
		encItemKeys := make([]string, 0)
		for _, itemKey := range panel {
			encrypted, err := unlynx.EncryptWithCothorityKey(itemKey)
			if err != nil {
				return nil, err
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanelsItemKeys = append(encPanelsItemKeys, encItemKeys)
	}
	return encPanelsItemKeys, nil
}
