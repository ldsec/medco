package survivalserver

const sql1 string = `
SELECT patient_num,start_date 
FROM observation_fact
WHERE concept_cd = $1 and patient_num = ANY($2::integer[])
`
const sql2 string = `
SELECT patient_num,end_date DATE_PART() FROM observation_fact
WHERE concept_cd = $3 and patient_num = ANY($2::integer[])
`
const sql3 string = `
SELECT DATE_PART(end_date::timestamp,start_date::timestamp) AS timepoint, COUNT(*) AS event_count
FROM (` + sql1 + `) AS x
INNER JOIN  (` + sql2 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const sql4 string = `
SELECT patient_num, MAX(end_date) AS end_date
WHERE patient_num patient_num = ANY($2::integer[]) AND patient_num NOT IN (SELECT patient_num FROM (` + sql2 + `))
GROUP BY patient_num
`

const sql5 string = `
SELECT DATE_PART(end_date::timestamp,start_date::timestamp) AS timepoint, COUNT(*) AS censoring_count
FROM (` + sql4 + `) AS x
INNER JOIN  (` + sql1 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const sql6 string = `
SELECT timepoint, event_count, censoring_count FROM (` + sql3 + `) AS x  FULL JOIN (` + sql5 + `) AS y
ON x.timepoint = y.timepoint
`
