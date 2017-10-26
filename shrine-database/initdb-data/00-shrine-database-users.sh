#!/bin/bash
set -e

mysql -p$ADMIN_PASSWORD -u root <<-EOSQL

    create database shrine_query_history;
    grant all privileges on shrine_query_history.* to shrine@% identified by '$DB_PASSWORD';
    create database stewardDB;
    grant all privileges on stewardDB.* to shrine@% identified by '$DB_PASSWORD';
    create database adapterAuditDB;
    grant all privileges on adapterAuditDB.* to shrine@% identified by '$DB_PASSWORD';
    create database qepAuditDB;
    grant all privileges on qepAuditDB.* to shrine@% identified by '$DB_PASSWORD';

EOSQL
