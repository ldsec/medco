package medcomodels

import (
	"time"

	"github.com/ldsec/medco/connector/restapi/models"
)

// Cohort holds cohort backend  reference
type Cohort struct {
	CohortID        int
	QueryID         int
	CohortName      string
	CreationDate    time.Time
	UpdateDate      time.Time
	QueryDefinition struct {
		Panels      []*models.Panel
		QueryTiming models.Timing
	}
}
