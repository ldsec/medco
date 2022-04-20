#!/bin/bash
set -Eeuo pipefail
# set up data in the database for end-to-end tests of the bioref feature
# Check this file for the loading of related i2b2 data /medco/build/package/i2b2/sql/86-e2etest-bioref.sh



psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL

        -- we create a bioref cohort containing all patients mentioned in the rows inserted into observation_fact


        INSERT INTO query_tools.explore_query_results(
        query_id, query_name, user_id, query_status, clear_result_set_size, clear_result_set, query_definition, i2b2_encrypted_patient_set_id, i2b2_non_encrypted_patient_set_id)
        VALUES (-42, 'BiorefE2ETest', 'test', 'completed', 9,'{1,2,3,4,5,6,7,8,9}', 'This field normally contains the selection panels', -1, -1)
        ON CONFLICT DO NOTHING;


        INSERT INTO query_tools.saved_cohorts(
        cohort_id, user_id, cohort_name, query_id, create_date, update_date, predefined, default_flag)
        VALUES (-42, 'test', 'testCohortBioref', -42, '2020-08-25 13:57:00', '2020-08-25 13:57:00', FALSE, FALSE)
        ON CONFLICT DO NOTHING;

        --it could be nice to create a second test cohort containing only a subset of the patients

EOSQL

