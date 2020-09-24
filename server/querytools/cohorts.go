package querytools

import "time"

// Cohort holds cohort backend  reference
type Cohort struct {
	CohortID     int
	QueryID      int
	CohortName   string
	CreationDate time.Time
	UpdateDate   time.Time
}
