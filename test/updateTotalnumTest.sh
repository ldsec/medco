#!/usr/bin/env bash


# This test script is to be executed from the root directory of the project
# This test script verifies that behaviour of the postgres update totalnum script is the one expected.
# This script does that by copying e2etest in a copy version of the table in order to keep the e2etest table as is.
# Then this script applied the update totalnum postgres script on this copied version of the table.
# When this is done we select all lines of interest in the ontology to verify that the c_totalnum attribute has been updated as expected.

MEDCO_DB_HOST="${1:-localhost}"
MEDCO_DB_PORT="${2:-5432}"
I2B2_DB_USER="${3:-i2b2}"
PGPASSWORD="${4:-i2b2}"
I2B2_DB_NAME="${5:-i2b2medcosrv0}"

expected_totalnum_patient="\e2etest\     |          4
\e2etest\1\   |          4
\e2etest\2\   |          4
\e2etest\3\   |          4
\modifiers\   |          4
\modifiers\1\ |          2
\modifiers\2\ |          2
\modifiers\3\ |          2"


expected_totalnum_observation="\e2etest\    |   17
\e2etest\1\ |   13
\e2etest\2\ |    12
\e2etest\3\ |    12
\modifiers\ |    10
\modifiers\1\   |  2
\modifiers\2\   |  2
\modifiers\3\   |  2"


table_copy=e2etest_testcpy

ontology_schema=medco_ont

export PSQL_PARAMS="-h ${MEDCO_DB_HOST} -p ${MEDCO_DB_PORT} -U ${I2B2_DB_USER}"
export I2B2_DB_USER

export PGPASSWORD

#updating the c_totalnum field in the database of the first local medco node
export I2B2_DB_NAME


#for each option p (patient) or o (observation) we test the psql totalnum updating script.
for option in p o; do
    python3 build/package/i2b2/sql/totalnum-update/generateUpdateTotalnumScript.py "${ontology_schema}" "${table_copy}" i2b2 i2b2demodata_i2b2 "${option}"


    psql $PSQL_PARAMS -d "${I2B2_DB_NAME}" <<-EOSQL
    DROP TABLE IF EXISTS ${ontology_schema}.${table_copy};
    CREATE TABLE ${ontology_schema}.${table_copy} AS TABLE ${ontology_schema}.e2etest;
EOSQL

    psql $PSQL_PARAMS -d "${I2B2_DB_NAME}" -f updateTotalnum.psql

    #verification phase we do some selects on e2etest data to check if totalnum was updated correctly.
    totalnums="$(
    psql $PSQL_PARAMS -d "${I2B2_DB_NAME}" <<-EOSQL
    SELECT c_fullname, c_totalnum FROM ${ontology_schema}.${table_copy}
        WHERE c_fullname LIKE '%modifiers\_\' ESCAPE '|' OR c_fullname LIKE '%e2etest\_\'  ESCAPE '|'
            OR c_fullname = '\e2etest\' OR c_fullname = '\modifiers\'
EOSQL
    )"

    trimmed_result=$(echo "$totalnums" | tr -d '\t -+')
    echo "Trimmed result: $trimmed_result"


    case "$option" in
    "p")
        expected_totalnum="${expected_totalnum_patient}"
        ;;
    "o")
        expected_totalnum="${expected_totalnum_observation}"
        ;;
    *)
        echo "there is a bug"
        exit 1
        #do nothing
        ;;
    esac

    expected_totalnum="$( echo "${expected_totalnum}" |  tr -d '\t ' )"

    #verifying that each line of the expected string is contained in the result
    while IFS= read -r expected_line; do
        if ! echo "${trimmed_result}" | grep "$expected_line" --quiet --fixed-strings; then
            echo "test failed line not contained in expected result:...$expected_line..."
            echo "expected result:${expected_totalnum}"
            exit 1
        fi
    done <<< "$expected_totalnum"

    #delete the table created for the sake of the test
    psql $PSQL_PARAMS -d "${I2B2_DB_NAME}" <<-EOSQL
    DROP TABLE ${ontology_schema}.${table_copy};
EOSQL

    echo "Test succeeded for option '${option}'."
    echo

done
