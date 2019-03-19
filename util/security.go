package util

import (
)

// --- PICSURE 2 internal authorization

// returns true if the token is valid
func ValidatePICSURE2InternalToken(token string) bool {
	return picsure2InternalToken() == token
}

// --- JWT-based MedCo user authentication
