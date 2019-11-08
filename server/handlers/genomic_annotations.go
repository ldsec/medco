package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
)

func MedCoGenomicAnnotationsGetValuesHandler(params genomic_annotations.GetValuesParams, principal *models.User) middleware.Responder {

	err := utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB")
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Impossible to connect to DB " + err.Error(),
		})
	}

	rows, err := utilserver.DBConnection.Query("SELECT annotation_value FROM genomic_annotations.? WHERE annotation_value = ~* ? LIMIT ?", params.Annotation, params.Value, *params.Limit)
	if err != nil {
		logrus.Error("Query execution error")
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query execution error " + err.Error(),
		})
	}
	defer rows.Close()

	var annotations []string
	var annotation string

	for rows.Next() {
		err := rows.Scan(&annotation)
		if err != nil {
			logrus.Error("Query result reading error")
			return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
				Message: "Query result reading error " + err.Error(),
			})
		}
		annotations = append(annotations, annotation)
	}

	if len(annotations) > 0 {
		return genomic_annotations.NewGetValuesOK().WithPayload(annotations)
	} else {
		return genomic_annotations.NewGetValuesNotFound()
	}
}

func MedCoGenomicAnnotationsGetVariantsHandler(params genomic_annotations.GetVariantsParams, principal *models.User) middleware.Responder {

	err := utilserver.DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB")
		return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
			Message: "Impossible to connect to DB " + err.Error(),
		})
	}

	rows, err := utilserver.DBConnection.Query("SELECT variant_id FROM genomic_annotations.genomic_annotations WHERE ? = ? AND  annotations = ~* ?", params.Annotation, params.Value, params.Zygosity)
	if err != nil {
		logrus.Error("Query execution error")
		return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
			Message: "Query execution error " + err.Error(),
		})
	}
	defer rows.Close()

	var variants []string
	var variant string

	for rows.Next() {
		err := rows.Scan(&variant)
		if err != nil {
			logrus.Error("Query result reading error")
			return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
				Message: "Query result reading error " + err.Error(),
			})
		}
		variants = append(variants, variant)
	}

	if len(variants) > 0 {
		return genomic_annotations.NewGetVariantsOK().WithPayload(variants)
	} else {
		return genomic_annotations.NewGetVariantsNotFound()
	}
}
