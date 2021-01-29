package survivalserver

import (
	"strings"
)

// getCodes
func getCode(path string) ([]string, error) {
	return
}

// getModifierCodes
func getModifierCodes(path string, appliedPath string) ([]string, error) {

}

// prepareLike prepare path for LIKE operator
func prepareLike(path string) string {
	if strings.HasSuffix(path, "%") {
		return path
	}
	if strings.HasSuffix(path, `\`) {
		return path + "%"
	}
	return path + `\%`
}
