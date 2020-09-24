package querytools

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_id = $2 AND query_status = 'completed');
`

const insertExploreResultInstance string = `
INSERT INTO query_tools.explore_query_results(user_id,query_name, query_status,query_definition)
VALUES ($1,$2,'running',$3)
RETURNING query_id
`

const updateExploreResultInstanceBoth string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' , i2b2_encrypted_patient_set_id=$4, i2b2_non_encrypted_patient_set_id=$5
WHERE query_id = $1 AND status = 'running'
`
const updateExploreResultInstanceOnlyClear string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' ,i2b2_non_encrypted_patient_set_id=$4
WHERE query_id = $1 AND query_status = 'running'
`
const updateExploreResultInstanceOnlyEncrypted string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' ,i2b2_encrypted_patient_set_id=$4
WHERE query_id = $1 AND query_status = 'running'
`

const updateErrorExploreQueryInstance string = `
UPDATE query_tools.explore_query_results
SET query_status='error'
WHERE query_id = $1 AND query_status = 'running'
`

const insertCohort string = `
INSERT INTO query_tools.saved_cohorts(user_id,query_id,cohort_name,create_date,update_date)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (user_id,cohort_name) DO UPDATE SET query_id = $2, update_date=$5
RETURNING cohort_id
`

const updateCohort string = `
UPDATE query_tools.saved_cohorts
SET query_id=$3, update_date= $4
WHERE cohort_id = $1 AND user_id = $2
`

const getCohorts string = `
SELECT cohort_id, query_id, cohort_name, create_date, update_date FROM query_tools.saved_cohorts
WHERE user_id = $1
`

const getDate string = `
SELECT update_date FROM query_tools.saved_cohorts
WHERE user_id =$1 and cohort_id=$2
`
