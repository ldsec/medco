package cohortscommon

import "time"

type Cohort struct {
	CohortId     int
	QueryId      int
	CohortName   string
	CreationDate time.Time
	UpdateDate   time.Time
}
