package directaccess

import (
	"fmt"
	"strings"

	survivalserver "github.com/ldsec/medco-connector/survival/server"
	utilserver "github.com/ldsec/medco-connector/util/server"
)

type DirectAccessTags func(map[string]string) (map[string]string, error)

func (tagHandler *DirectAccessTags) GetTagIDs(tags map[string]string) (map[string]string, error) {
	return (*tagHandler)(tags)
}

const (
	conceptTable string = "concept_dimension"
	pathColumn   string = "concept_path"
	tagQuery     string = `SELECT concept_path,concept_cd FROM ` + schema + `.` + conceptTable
)

type tagRecipiens struct {
	path        string
	conceptCode string
}

func getTagIDs(tags map[string]string) (tagIDs map[string]string, err error) {

	paths := buildParameters(tags)
	psqlQuery := tagQuery + ` WHERE concept_path IN (` + paths + `)`
	rows, err := DirectAccessDB.Query(psqlQuery)
	err = survivalserver.NiceError(err)
	if err != nil {
		return
	}
	tagIDs = make(map[string]string, len(tags))
	numberOfRows := 0

	for rows.Next() {
		recipiens := &tagRecipiens{}
		err = rows.Scan(&(recipiens.path), &(recipiens.conceptCode))

		if err != nil {
			return
		}
		tag := strings.Replace(recipiens.path, "\\medco\\tagged\\concept\\", "", 1)
		tag = strings.Replace(tag, "\\", "", 1)
		encTimeCode, ok := tags[tag]
		if !ok {
			err = fmt.Errorf("tag  %s  not found, node index %d", tag, utilserver.MedCoNodeIdx)
			return
		}

		tagIDs[recipiens.conceptCode] = encTimeCode
		numberOfRows++

	}
	if numberOfRows == 0 {
		err = fmt.Errorf("From node %d, Unable to find any of the tag in the data base %s", utilserver.MedCoNodeIdx, DBName)
	}

	err = rows.Close()

	return

}

func buildParameters(tags map[string]string) string {
	paths := make([]string, len(tags))
	pos := 0
	for tag := range tags {
		paths[pos] = `'\medco\tagged\concept\` + tag + `\'`
		pos++
	}

	return strings.Join(paths, ",")
}
