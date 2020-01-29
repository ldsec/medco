#!/usr/bin/env bash

DB_NUMBER=9

#check if the number of connections to the postgresql server is at least as much as DB_NUMBER to detect that the startup phase has begun
while true; do
  echo "Starting MedCo up..."
  DB_CONNECTIONS_NUMBER=$(PGPASSWORD=postgres psql -v ON_ERROR_STOP=1 -h localhost -p 5432 -U postgres -d postgres -t -c "SELECT count(*) FROM pg_stat_activity WHERE usename!='postgres';")
  if [ $DB_CONNECTIONS_NUMBER -gt $DB_NUMBER ]; then
    break
  fi
  sleep 15
done

counter=0

#periodically check the number of connections to the postgresql server. If no connections are detected for a minute, then we assume that the startup phase has been completed
while true; do
  echo "Waiting for MedCo startup to complete..."
  DB_CONNECTIONS_NUMBER=$(PGPASSWORD=postgres psql -v ON_ERROR_STOP=1 -h localhost -p 5432 -U postgres -d postgres -t -c "SELECT count(*) FROM pg_stat_activity WHERE (state='active' OR state='idle in transaction') AND usename!='postgres';")
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