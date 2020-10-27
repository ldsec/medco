package unlynx

// EncryptMatrix encrypts matrix entriess with the cothority key
func EncryptMatrix(itemKeys [][]int64) ([][]string, error) {
	encPanelsItemKeys := make([][]string, 0)
	for _, keyList := range itemKeys {
		encItemKeys := make([]string, 0)
		for _, itemKey := range keyList {
			encrypted, err := EncryptWithCothorityKey(itemKey)
			if err != nil {
				return nil, err
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanelsItemKeys = append(encPanelsItemKeys, encItemKeys)
	}
	return encPanelsItemKeys, nil
}
