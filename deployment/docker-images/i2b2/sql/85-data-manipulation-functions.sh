#!/bin/bash
set -Eeuo pipefail
# add the functions to manipulate the i2b2 data

for f in "${I2B2_SQL_DIR}/85-data-manipulation-functions/"/*; do
    psql ${PSQL_PARAMS} -d "${I2B2_DB_NAME}" -f "$f"
done
