package survivalserver

const sql1 string = `
SELECT patient_num,start_date 
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = $1 and modifier_cd = $2 and patient_num = ANY($3::integer[])
`
const sql2 string = `
SELECT patient_num,end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = $4 and modifier_cd = $5 and patient_num = ANY($3::integer[])
`
const sql3 string = `
SELECT DATE_PART('day',end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS event_count
FROM (` + sql1 + `) AS x
INNER JOIN  (` + sql2 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const sql4 string = `
SELECT patient_num, MAX(end_date) AS end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE patient_num = ANY($3::integer[]) AND patient_num NOT IN (SELECT patient_num FROM (` + sql2 + `) AS patients_with_events)
GROUP BY patient_num
`

const sql5 string = `
SELECT DATE_PART('day', end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS censoring_count
FROM (` + sql4 + `) AS x
INNER JOIN  (` + sql1 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const sql6 string = `
SELECT COALESCE(xx.timepoint,yy.timepoint) AS timepoint , COALESCE(event_count,0) AS event_count, COALESCE(censoring_count,0) AS censoring_count FROM (` + sql3 + `) AS xx  FULL JOIN (` + sql5 + `) AS yy
ON xx.timepoint = yy.timepoint
`
