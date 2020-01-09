package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
	"strings"
)

// MedCoGenomicAnnotationsGetValuesHandler handles /medco/genomic-annotations/{annotation} API endpoint
func MedCoGenomicAnnotationsGetValuesHandler(params genomic_annotations.GetValuesParams, principal *models.User) middleware.Responder {

	err := utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Impossible to connect to DB " + err.Error(),
		})
	}

	var query string

	query, err = buildGetValuesQuery(params)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetVariantsNotFound()
	}

	//escaping * characters
	rows, err := utilserver.DBConnection.Query(query, params.Annotation, strings.ReplaceAll(params.Value, "*", "\\*"), *params.Limit)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query execution error: " + err.Error(),
		})
	}
	defer rows.Close()

	var annotations []string
	var annotation string

	for rows.Next() {
		err := rows.Scan(&annotation)
		if err != nil {
			logrus.Error("Query result reading error: " + err.Error())
			return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
				Message: "Query result reading error: " + err.Error(),
			})
		}
		annotations = append(annotations, annotation)
	}

	return genomic_annotations.NewGetValuesOK().WithPayload(annotations)

}

// MedCoGenomicAnnotationsGetVariantsHandler handles /medco/genomic-annotations/{annotation}/{value} API endpoint
func MedCoGenomicAnnotationsGetVariantsHandler(params genomic_annotations.GetVariantsParams, principal *models.User) middleware.Responder {

	err := utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB: " + err.Error())
		return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
			Message: "Impossible to connect to DB: " + err.Error(),
		})
	}
	zygosityStr := ""
	if len(params.Zygosity) > 0 {
		zygosityStr = params.Zygosity[0]

		for i := 1; i < len(params.Zygosity); i++ {
			zygosityStr += "|" + params.Zygosity[i]
		}
	}

	var query string

	query, err = buildGetVariantsQuery(params)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetVariantsNotFound()
	}

	rows, err := utilserver.DBConnection.Query(query, params.Annotation, params.Value, zygosityStr, *params.Encrypted)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query execution error: " + err.Error(),
		})
	}
	defer rows.Close()

	var variants []string
	var variant string

	for rows.Next() {
		err := rows.Scan(&variant)
		if err != nil {
			logrus.Error("Query result reading error: " + err.Error())
			return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
				Message: "Query result reading error: " + err.Error(),
			})
		}
		variants = append(variants, variant)
	}

	if len(variants) > 0 {
		return genomic_annotations.NewGetVariantsOK().WithPayload(variants)
	}
	return genomic_annotations.NewGetVariantsNotFound()

}

func buildGetValuesQuery(params genomic_annotations.GetValuesParams) (string, error) {

	if contains(utilserver.GenomicAnnotationTypes, params.Annotation) {
		return "SELECT \"ga_getvalues\"($1,$2,$3)", nil
	}
	return "", errors.New("Requested invalid annotation type: " + params.Annotation)

}

func buildGetVariantsQuery(params genomic_annotations.GetVariantsParams) (string, error) {

	if contains(utilserver.GenomicAnnotationTypes, params.Annotation) {
		return "SELECT \"ga_getvariants\"($1,$2,$3,$4)", nil
	}
	return "", errors.New("Requested invalid annotation type: " + params.Annotation)

}

func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}
