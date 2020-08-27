package querytools

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_id = $2 AND query_status = 'completed');
`

const insertExploreResultInstance string = `
INSERT INTO query_tools.explre_query_results(user_id,query_name, query_status,query_definition)
VALUES ($1,$2,'running',$3)
RETURNING query_id
`

const updateExploreResultInstance string = `
UPDATE medco_query_tools.qt_query_instance
SET (clear_result_set_size, clear_result_set,query_status,i2b2_encrypted_patient_set_id,i2b2_non_encrypted_patient_set_id)=$2, $3, 'completed', $4, $5
WHERE query_instance_id = $1
`
const updateErrorExploreQueryInstance string = `
UPDATE medco_query_tools.qt_query_instance
SET (query_status,i2b2_encrypted_patient_set_id,i2b2_non_encrypted_patient_set_id)='error',$2,$3
WHERE query_instance_id = $1
`

const insertCohort string = `
INSERT INTO query_tools.saved_cohorts(cohort_id,user_id,query_id,cohort_name,create_date,update_date)
VALUES ($1,$2,$3,$4)
RETURNING cohort_id
`

const updateCohort string = `
UPDATE query_tools.saved_cohorts
SET (update_date) $3
WHERE cohort_id = $1 AND user_id = $2
`

const getCohorts string = `
SELECT cohort_id, cohort_name, create_date, update_date FROM query_tools.saved_cohorts
WHERE user_id = $1
`
