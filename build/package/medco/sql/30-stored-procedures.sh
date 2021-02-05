#!/bin/bash
set -Eeuo pipefail
# add the functions for query tools and saved cohorts features

for f in "${MC_SQL_DIR}/30-stored-procedures/"/*; do
    psql ${PSQL_PARAMS} -d "${MC_DB_NAME}" -f "$f"
done
