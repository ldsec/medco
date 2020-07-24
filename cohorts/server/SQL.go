package cohortsserver

const getCohorts string = `
SELECT cohort_name, result_instance_id, create_date, update_date
FROM medco_query_tools.cohorts
WHERE user_id = ?;
`

const getDates string = `
SELECT create_date, update_date
FROM medco_query_tools.cohorts
WHERE user_id = ? AND cohort_name = ?; 
`

const insertCohorts string = `
INSERT INTO cohorts (user_id,cohort_name,result_instance_id, create_date, update_date )
VALUES (?, ?, ?, ?, ?)
ON CONFLICT ON cohorts_user_id_cohort_name_key 
DO UPDATE SET result_instance_id = EXCLUDED.result_instance_id , update_date=EXCLUDED.update_date;
`
