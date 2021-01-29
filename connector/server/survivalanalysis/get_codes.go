package survivalserver

import (
	"fmt"
	"strings"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

// getCodes
func getCodes(path string) ([]string, error) {

	preparedPath := prepareLike(path)
	description := fmt.Sprintf("getCodes (path: %s), SQL: %s", preparedPath, codes)
	logrus.Debugf("running: %s", description)
	rows, err := utilserver.I2B2DBConnection.Query(codes, preparedPath)
	if err != nil {
		err = fmt.Errorf("while selecting concept codes: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}

	resString := new(string)
	res := make([]string, 0)
	for rows.Next() {

		err = rows.Scan(resString)
		if err != nil {
			err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
			return nil, err
		}

		res = append(res, *resString)
	}
	logrus.Tracef("concept codes are %v, DB operation: %s", res, description)
	err = rows.Close()
	if err != nil {
		err = fmt.Errorf("while closing SQL record stream: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}
	logrus.Debugf("successfully retrieved %d concept codes, DB operation: %s", len(res), description)

	return res, nil

}

// getModifierCodes
func getModifierCodes(path string, appliedPath string) ([]string, error) {
	preparedPath := prepareLike(path)
	preparedAppliedPath := prepareEqual(appliedPath)
	description := fmt.Sprintf("getModifierCodes (path: %s, appliedPath: %s), SQL: %s", preparedPath, preparedAppliedPath, modifierCodes)
	logrus.Debugf("running: %s", description)
	rows, err := utilserver.I2B2DBConnection.Query(modifierCodes, preparedPath, preparedAppliedPath)
	if err != nil {
		err = fmt.Errorf("while selecting modifier codes: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}

	resString := new(string)
	res := make([]string, 0)
	for rows.Next() {

		err = rows.Scan(resString)
		if err != nil {
			err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
			return nil, err
		}

		res = append(res, *resString)
	}
	logrus.Tracef("modifier codes are %v, DB operation: %s", res, description)

	err = rows.Close()
	if err != nil {
		err = fmt.Errorf("while closing SQL record stream: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}

	logrus.Debugf("successfully retrieved %d modifier codes, DB operation: %s", len(res), description)
	return res, nil
}

// prepareLike prepare path for LIKE operator
func prepareLike(pathURI string) string {
	path := strings.Replace(pathURI, "/", `\\`, -1)
	if strings.HasSuffix(path, "%") {
		return path
	}
	if strings.HasSuffix(path, `\`) {
		return path + "%"
	}
	return path + `\\%`
}

// prepareEqual prepare path for LIKE operator
func prepareEqual(pathURI string) string {
	return strings.Replace(pathURI, "/", `\`, -1)
}

const codes = `
SELECT c_basecode
	FROM medco_ont.e2etest
	WHERE (c_basecode IS NOT NULL AND c_basecode  != ''
			AND c_facttablecolumn = 'concept_cd'
		  AND c_fullname LIKE $1);
`

const modifierCodes = `
SELECT c_basecode
	FROM medco_ont.e2etest
	WHERE (c_basecode IS NOT NULL AND c_basecode  != ''
			AND c_facttablecolumn = 'modifier_cd'
		  AND c_fullname LIKE $1 AND m_applied_path = $2);
`
