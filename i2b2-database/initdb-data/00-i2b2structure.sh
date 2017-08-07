#!/bin/bash
set -e

# TODO: password hardcoded to be modified (and passed by secrets of docker)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE $I2B2_DOMAIN_NAME
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    create schema i2b2demodata;
    create schema i2b2hive;
    create schema i2b2imdata;
    create schema i2b2metadata;
    create schema i2b2pm;
    create schema i2b2workdata;

    create role i2b2demodata login password 'demouser';
    create role i2b2hive login password 'demouser';
    create role i2b2imdata login password 'demouser';
    create role i2b2metadata login password 'demouser';
    create role i2b2pm login password 'demouser';
    create role i2b2workdata login password 'demouser';

    grant all on schema i2b2demodata to i2b2demodata;
    grant all on schema i2b2hive to i2b2hive;
    grant all on schema i2b2imdata to i2b2imdata;
    grant all on schema i2b2metadata to i2b2metadata;
    grant all on schema i2b2pm to i2b2pm;
    grant all on schema i2b2workdata to i2b2workdata;

    grant all privileges on all tables in schema i2b2demodata to i2b2demodata;
    grant all privileges on all tables in schema i2b2hive to i2b2hive;
    grant all privileges on all tables in schema i2b2imdata to i2b2imdata;
    grant all privileges on all tables in schema i2b2metadata to i2b2metadata;
    grant all privileges on all tables in schema i2b2pm to i2b2pm;
    grant all privileges on all tables in schema i2b2workdata to i2b2workdata;
EOSQL

