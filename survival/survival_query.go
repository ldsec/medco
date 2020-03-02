package survival

import (
	"time"

	survivalserver "github.com/ldsec/medco-connector/survival/server"
	"github.com/ldsec/medco-connector/survival/server/directaccess"
)

type Query struct {
	ID            string
	UserPublicKey string
	PatientSetID  string
	TimeCodes     []string

	//TODO also hide that
	Result *struct {
		Timers    map[string]time.Duration
		EncEvents map[string][2]string
	}

	spin *survivalserver.Spin
}

func NewQuery(qID, pubKey, patientSetID string, timeCodes []string) *Query {
	return &Query{ID: qID, UserPublicKey: pubKey, PatientSetID: patientSetID, TimeCodes: timeCodes, spin: survivalserver.NewSpin()}
}

func (q *Query) GetID() string {
	return q.ID
}

func (q *Query) GetUserPublicKey() string {
	return q.UserPublicKey
}

func (q *Query) GetPatientSetID() string {
	return q.PatientSetID
}

func (q *Query) GetTimeCodes() []string {
	return q.TimeCodes
}

//TODO sync map this
func (q *Query) SetResultMap(resultMap map[string][2]string) error {
	q.spin.Lock()
	defer q.spin.Unlock()
	if q.Result == nil {
		q.Result = new(struct {
			Timers    map[string]time.Duration
			EncEvents map[string][2]string
		})
	}
	q.Result.EncEvents = resultMap
	return nil

}

func (q *Query) Execute() error {
	err := directaccess.QueryTimePoints(q, 1)

	return err
}
