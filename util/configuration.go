package util

import "os"

// public ------------

// i2b2 connection ---

// get the URL of the i2b2 CRC cell this connector is using
func I2b2CRCCellURL() string {
	return os.Getenv("I2B2_CRC_URL")
}

// get the URL of the i2b2 ONT cell this connector is using
func I2b2ONTCellURL() string {
	return os.Getenv("I2B2_ONT_URL")
}

// get the login domain for i2b2
func I2b2LoginDomain() string {
	return os.Getenv("I2B2_LOGIN_DOMAIN")
}
// get the login project for i2b2
func I2b2LoginProject() string {
	return os.Getenv("I2B2_LOGIN_PROJECT")
}
// get the login user for i2b2
func I2b2LoginUser() string {
	return os.Getenv("I2B2_LOGIN_USER")
}
// get the login password
func I2b2LoginPassword() string {
	return os.Getenv("I2B2_LOGIN_PASSWORD")
}

// get the timeout in seconds for communications with i2b2
func I2b2TimeoutSeconds() int {
	return 180
}

// private ------------

// get the token (shared secret) used for internal PICSURE 2 authorization
func picsure2InternalToken() string {
	return os.Getenv("PICSURE2_INTERNAL_TOKEN")
}

