#!/usr/bin/env bash
set -Euo pipefail

DB_NUMBER=${1:-9}
DB_HOST=${2:-localhost}
DB_PORT=${3:-5432}
DB_USER=${4:-postgres}
DB_PASSWORD=${5:-postgres}
DB_NAME=${6:-postgres}

#check if the number of connections to the postgresql server is at least as much as DB_NUMBER to detect that the startup phase has begun
while true; do
  echo "Starting MedCo up..."
  DB_CONNECTIONS_NUMBER=$(PGPASSWORD=$DB_PASSWORD psql -v ON_ERROR_STOP=1 -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM pg_stat_activity WHERE usename!='$DB_USER';")
  if [ $DB_CONNECTIONS_NUMBER -gt $DB_NUMBER ]; then
    break
  fi
  sleep 15
done

counter=0

#periodically check the number of connections to the postgresql server. If no connections are detected for a minute, then we assume that the startup phase has been completed
while true; do
  echo "Waiting for MedCo startup to complete..."
  DB_CONNECTIONS_NUMBER=$(PGPASSWORD=$DB_PASSWORD psql -v ON_ERROR_STOP=1 -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM pg_stat_activity WHERE (state='active' OR state='idle in transaction') AND usename!='$DB_USER';")
  if [ $DB_CONNECTIONS_NUMBER -eq 0 ]; then
    if [ $counter -eq 3 ]; then
      break
    fi
    ((counter++))
  else
    counter=0
  fi
  sleep 15
done

echo "---MedCo is up and running!---"