#!/usr/bin/env bash
set -Eeuo pipefail
# example: bash load-predefined-filters.sh localhost medcoconnectorsrv0 ../test/test_predefined/predefined_cohorts.txt
# loads the the test file predefined_cohorts.txt in the first medco node of the local dev deployment, it will fails if run more than once
# The file that contains the predefined filters must holds one filter per line, with its name and definition
# separated by a comma ",", without whitespace before or after like the following example
# filterName,{"anyPanelDefnintionObject":"anyValue"}
#

DB_HOST=$1
MEDCOCONNECTORDB=$2
FILE_PATH=$3


QUERY_NUMBER=0
QUERY_DATE=$(/usr/bin/date +%Y%m%d%H%M%S)


while IFS="" read -r p || [ -n "$p" ]
do

QUERY_NAME=predefined_${QUERY_DATE}_$((++QUERY_NUMBER))

# these sed functions separate the name from the panels
DEFINITION=$(sed 's/^[^,]*,//g' <<< $p)
FILTER_NAME=$(sed 's/,.*$//g' <<< $p)

PGPASSWORD=medcoconnector psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "medcoconnector" -p 5432 -d "${MEDCOCONNECTORDB}" <<-EOSQL
BEGIN;
INSERT INTO query_tools.explore_query_results(
	query_name, user_id, query_status, query_definition)
	VALUES ('${QUERY_NAME}', 'allUsers', 'predefined', '${DEFINITION}');

INSERT INTO query_tools.saved_cohorts(
	user_id, cohort_name, query_id,  predefined)
	VALUES ('allUsers', '${FILTER_NAME}',
    (SELECT query_id FROM query_tools.explore_query_results WHERE query_name = '${QUERY_NAME}'),
  TRUE);
COMMIT;
EOSQL

done < $FILE_PATH