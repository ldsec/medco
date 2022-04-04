#!/bin/bash
set -Eeuo pipefail

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  -- create roles and give them rights to create database
  CREATE ROLE i2b2 LOGIN PASSWORD '$DB_ALL_USERS_PASSWORD';
  CREATE ROLE keycloak LOGIN PASSWORD '$DB_ALL_USERS_PASSWORD';
  CREATE ROLE medcoconnector LOGIN PASSWORD '$DB_ALL_USERS_PASSWORD';
  ALTER USER i2b2 CREATEDB;
  ALTER USER keycloak CREATEDB;
  ALTER USER medcoconnector CREATEDB;

  -- create database for keycloak
  CREATE DATABASE keycloak;
  ALTER DATABASE keycloak OWNER TO keycloak;
EOSQL
