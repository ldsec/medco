package querytoolsserver

import (
	"fmt"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

// CohortsPatientList holds the parameters structure to retriev the patient list of a saved cohort
type CohortsPatientList struct {
	CohortName    string
	ID            string
	User          *models.User
	UserPublicKey string

	Result struct {
		EncPatientList []string
		Timers         medcomodels.Timers
	}
}

// NewCohortsPatientList returns a new pointer on a CohortsPatientList
func NewCohortsPatientList(ID, UserPublicKey, cohortName string, user *models.User) (res *CohortsPatientList) {
	res = &CohortsPatientList{
		ID:            ID,
		User:          user,
		CohortName:    cohortName,
		UserPublicKey: UserPublicKey,
	}

	res.Result.EncPatientList = make([]string, 0)
	res.Result.Timers = medcomodels.NewTimers()

	return res
}

// Execute retrieves patient numbers and encrypt them
func (pl *CohortsPatientList) Execute() error {
	// get patient list number
	timer := time.Now()
	pIDs, err := GetPatientList(pl.User.ID, pl.CohortName)
	if err != nil {
		err = fmt.Errorf("while requesting patient list for cohort patient list retrieval %s: %s", pl.ID, err.Error())
		logrus.Error(err.Error())
		return err
	}
	pl.Result.Timers.AddTimers("get-patient-numbers", timer, nil)

	// encryption. It is assumed there is no dummy patients.
	encryptedIDs := make([]string, len(pIDs))
	timer = time.Now()
	for i, pID := range pIDs {
		encryptedIDs[i], err = unlynx.EncryptWithCothorityKey(pID)
		if err != nil {
			// error itself not returned to client to avoid leaking information
			err = fmt.Errorf("during encryption of patient number for cohort patient list retrieval %s", pl.ID)
			logrus.Error(err.Error())
			return err
		}
	}
	pl.Result.Timers.AddTimers("encryption", timer, nil)

	// Keyswitch
	timer = time.Now()
	encPatientList, ksTimers, err := unlynx.KeySwitchValues(pl.ID, encryptedIDs, pl.UserPublicKey)
	if err != nil {
		err = fmt.Errorf("during encryption of patient number for cohort patient list retrieval %s: %s", pl.ID, err.Error())
		logrus.Error(err.Error())
		return err
	}
	pl.Result.Timers.AddTimers("key-switch", timer, ksTimers)
	pl.Result.EncPatientList = encPatientList
	return nil
}
