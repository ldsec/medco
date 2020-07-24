package querytools

const getPatientList string = `
SELECT patient_num FROM medco_query_tools.patient_set_collection
WHERE result_instance_id  = $1
`

const insertQueryInstance string = `
INSERT INTO medco_query_tools.qt_query_instance(user_id, group_id, batchmode, start_date)
VALUES ($1,$2,$3,$4)
RETURNING query_instance_id
`

const updateQueryInstance string = `
UPDATE medco_query_tools.qt_query_instance
SET (clear_query_instance_id, genomic_query_instance_id,end_date,batch_mode,status_type,message)=$2, $3, $4, $5, $6
WHERE query_instance_id = $1
`

const insertResultInstance string = `
INSERT INTO medco_query_tools.qt_result_instance(result_type, query_instance_id, start_date, status_type_id)
VALUES ($1,$2,$3,$4)
RETURNING result_instance_id
`

const updateResultInstance string = `
UPDATE medco_query_tools.qt_result_instance
SET (set_size, end_date, status_type_id, message, real_set_size, obfusctation_method) $2, $3, $4, $5, $6, $7
WHERE result_instance_id = $1
`

/*
 * Easiest way to insert patient in the patient set Collection, not the optimal
 */
const insertPatientSetCollection string = `
INSER INTO medco_query_tools.qt_patient_set_collection(result_instance_id, set_index, patient_num)
VALUES ($1, $2, $3)
RETURNING patient_set_coll_id
`
