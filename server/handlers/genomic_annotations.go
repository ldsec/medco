package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	"github.com/ldsec/medco-connector/server/genomicannotations"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"
)

// MedCoGenomicAnnotationsGetValuesHandler handles /medco/genomic-annotations/{annotation} API endpoint
func MedCoGenomicAnnotationsGetValuesHandler(params genomic_annotations.GetValuesParams, principal *models.User) middleware.Responder {

	err := utilserver.GaDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Impossible to connect to DB " + err.Error(),
		})
	}

	query, err := genomicannotations.BuildGetValuesQuery(params)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetVariantsNotFound()
	}

	rows, err := genomicannotations.ExecuteGetValuesQuery(query, params)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query execution error: " + err.Error(),
		})
	}
	defer rows.Close()

	annotations, err := genomicannotations.BuildGetValuesQueryResponse(rows)
	if err != nil {
		logrus.Error("Query result reading error: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query result reading error: " + err.Error(),
		})
	}

	return genomic_annotations.NewGetValuesOK().WithPayload(annotations)

}

// MedCoGenomicAnnotationsGetVariantsHandler handles /medco/genomic-annotations/{annotation}/{value} API endpoint
func MedCoGenomicAnnotationsGetVariantsHandler(params genomic_annotations.GetVariantsParams, principal *models.User) middleware.Responder {

	err := utilserver.GaDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to DB: " + err.Error())
		return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
			Message: "Impossible to connect to DB: " + err.Error(),
		})
	}

	query, err := genomicannotations.BuildGetVariantsQuery(params)

	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetVariantsNotFound()
	}

	rows, err := genomicannotations.ExecuteGetVariantsQuery(query, params)
	if err != nil {
		logrus.Error("Query execution error: " + err.Error())
		return genomic_annotations.NewGetValuesDefault(500).WithPayload(&genomic_annotations.GetValuesDefaultBody{
			Message: "Query execution error: " + err.Error(),
		})
	}
	defer rows.Close()

	variants, err := genomicannotations.BuildGetVariantsQueryResponse(rows)
	if err != nil {
		logrus.Error("Query result reading error: " + err.Error())
		return genomic_annotations.NewGetVariantsDefault(500).WithPayload(&genomic_annotations.GetVariantsDefaultBody{
			Message: "Query result reading error: " + err.Error(),
		})
	}

	if len(variants) > 0 {
		return genomic_annotations.NewGetVariantsOK().WithPayload(variants)
	}
	return genomic_annotations.NewGetVariantsNotFound()

}
