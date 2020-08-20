package survivalserver

import (
	"database/sql"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

/**
*  TestSql is aimed at unit tests run at development time. It is not expected to succeed in CI or other automated run env
 */
func TestSql(t *testing.T) {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=i2b2medcosrv0 sslmode=disable")
	logrus.Info(err)
	if err != nil {
		t.Error("Error at creating database", err)
	}
	err = db.Ping()
	if err != nil {
		t.Skip("Unable to connect database for testing", err)
	}

	res, err := db.Exec("CREATE SCHEMA survival_test")
	logrus.Info(res, err)
	path, err := exec.LookPath("./load_test_data.sh")
	if err != nil {
		t.Error("Error when finding ./load_test_data.sh", err)
	}
	logrus.Info("Found " + path)
	command := exec.Command("./load_test_data.sh")
	err = command.Start()
	if err != nil {
		t.Error("Error at starting command", err)
	}
	logrus.Info("Started data loading")
	err = command.Wait()
	if err != nil {
		t.Error("Command terminated with error", err)
	}
	count := new(int)
	err = db.QueryRow(`SELECT COUNT(*) FROM survival_test.observation_fact`).Scan(count)
	if err != nil {
		t.Error(err)
	}
	logrus.Infof("Found %d rows", *count)
	if *count != 393 {
		t.Errorf("Rows number in test observation_fact failed. Expected 393, got %d", *count)
	}
	patient := new(int)
	date := new(string)
	rows, err := db.Query(testSql1, "NDC:00015345620", "{483573, 483574, 483575}")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(patient, date)
		logrus.Info(*patient, *date)

	}
	timePoint := new(int)
	rows, err = db.Query(testSql3, "NDC:00015345620", "{483573, 483574, 483575}", "DEM|DATE:death")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(timePoint, count)
		logrus.Info(*timePoint, *count)

	}

	rows, err = db.Query(testSql5, "NDC:00015345620", "{483573, 483574, 483575}", "DEM|DATE:death")
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		rows.Scan(timePoint, count)
		logrus.Info(*timePoint, *count)
	}

	eventCount := count
	censoringCount := new(int)

	rows, err = db.Query(testSql6, "NDC:00015345620", "{483573, 483574, 483575}", "DEM|DATE:death")
	if err != nil {
		t.Error(err)
	}

	for rows.Next() {
		rows.Scan(timePoint, eventCount, censoringCount)
		logrus.Info(*timePoint, *eventCount, *censoringCount)
	}

	res, err = db.Exec("DROP SCHEMA survival_test CASCADE")
	logrus.Info(res, err)

	t.Skip("Unable to connect database", err)
}

const testSql1 string = `
SELECT patient_num,start_date 
FROM survival_test.observation_fact
WHERE concept_cd = $1 and patient_num = ANY($2::integer[])
`
const testSql2 string = `
SELECT patient_num,end_date
FROM survival_test.observation_fact
WHERE concept_cd = $3 and patient_num = ANY($2::integer[])
`
const testSql3 string = `
SELECT DATE_PART('day',end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS event_count
FROM (` + testSql1 + `) AS x
INNER JOIN  (` + testSql2 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const testSql4 string = `
SELECT patient_num, MAX(end_date) AS end_date
FROM survival_test.observation_fact
WHERE patient_num = ANY($2::integer[]) AND patient_num NOT IN (SELECT patient_num FROM (` + testSql2 + `) AS patients_with_events)
GROUP BY patient_num
`

const testSql5 string = `
SELECT DATE_PART('day', end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS censoring_count
FROM (` + testSql4 + `) AS x
INNER JOIN  (` + testSql1 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const testSql6 string = `
SELECT COALESCE(x.timepoint,y.timepoint) AS timepoint , COALESCE(event_count,0) AS event_count, COALESCE(censoring_count,0) AS censoring_count FROM (` + testSql3 + `) AS x  FULL JOIN (` + testSql5 + `) AS y
ON x.timepoint = y.timepoint
`
