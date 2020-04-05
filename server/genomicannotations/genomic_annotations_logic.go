package genomicannotations

import (
	"database/sql"
	"errors"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
	"strings"
)

// BuildGetValuesQuery builds the get-values query
func BuildGetValuesQuery(params genomic_annotations.GetValuesParams) (string, error) {

	if ok, err := annotationExists(params.Annotation); err != nil {
		logrus.Error("annotation exists check failed: ", err)
		return "", err
	} else if !ok {
		return "", errors.New("Requested invalid annotation type: " + params.Annotation)
	}

	return "SELECT genomic_annotations.ga_getvalues($1,$2,$3,$4)", nil
}

// ExecuteGetValuesQuery executes the get-values query
func ExecuteGetValuesQuery(query string, params genomic_annotations.GetValuesParams) (*sql.Rows, error) {

	//escaping * characters
	return utilserver.GaDBConnection.Query(query, params.Annotation, strings.ReplaceAll(params.Value, "*", "\\*"), *params.Limit, false)
}

// BuildGetValuesQueryResponse builds the response to the get-values query
func BuildGetValuesQueryResponse(rows *sql.Rows) (annotations []string, err error) {

	var annotation string

	for rows.Next() {
		err = rows.Scan(&annotation)
		if err != nil {
			return
		}
		annotations = append(annotations, annotation)
	}

	return
}

// BuildGetVariantsQuery builds the get-variants query
func BuildGetVariantsQuery(params genomic_annotations.GetVariantsParams) (string, error) {

	if ok, err := annotationExists(params.Annotation); err != nil {
		logrus.Error("annotation exists check failed: ", err)
		return "", err
	} else if !ok {
		return "", errors.New("Requested invalid annotation type: " + params.Annotation)
	}

	return "SELECT genomic_annotations.ga_getvariants($1,$2,$3,$4,$5)", nil
}

// ExecuteGetVariantsQuery executes the get-variants query
func ExecuteGetVariantsQuery(query string, params genomic_annotations.GetVariantsParams) (*sql.Rows, error) {

	zygosityStr := ""
	if len(params.Zygosity) > 0 {
		zygosityStr = params.Zygosity[0]
		for i := 1; i < len(params.Zygosity); i++ {
			zygosityStr += "|" + params.Zygosity[i]
		}
	}

	return utilserver.GaDBConnection.Query(query, params.Annotation, params.Value, zygosityStr, *params.Encrypted, false)
}

// BuildGetVariantsQueryResponse builds the response to the get-variants query
func BuildGetVariantsQueryResponse(rows *sql.Rows) (variants []string, err error) {

	var variant string

	for rows.Next() {
		err = rows.Scan(&variant)
		if err != nil {
			return
		}
		variants = append(variants, variant)
	}

	return
}

// annotationExists checks in the database if the annotation exists
func annotationExists(annotationName string) (exists bool, err error) {
	err = utilserver.GaDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB: " + err.Error())
		return
	}

	res, err := utilserver.GaDBConnection.Query("SELECT genomic_annotations.ga_annotationexists($1)", annotationName)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return
	}
	defer res.Close()

	if !res.Next() {
		err = errors.New("No result available for annotationexists check")
		logrus.Error(err)
		return
	}

	err = res.Scan(&exists)
	if err != nil {
		logrus.Error("Query result reading error: " + err.Error())
		return
	}
	return
}
