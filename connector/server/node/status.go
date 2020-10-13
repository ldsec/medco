package node

// CheckStatus checks the status of the MedCo node
func CheckStatus() (message string, status bool) {

	okI2b2 := statusI2b2()
	message += "i2b2: "
	if okI2b2 {
		message += "OK\n"
	} else {
		message += "error\n"
	}

	okUnlynx := statusUnlynx()
	message += "Unlynx: "
	if okUnlynx {
		message += "OK\n"
	} else {
		message += "error\n"
	}

	okGenomicAnnotations := statusGenomicAnnotations()
	message += "Genomic annotations: "
	if okGenomicAnnotations {
		message += "OK\n"
	} else {
		message += "error\n"
	}

	return message, okI2b2 && okUnlynx && okGenomicAnnotations

}
