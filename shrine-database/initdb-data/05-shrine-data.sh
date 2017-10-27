#!/bin/bash
set -e

# todo: some tables are duplicated in several databases
#mysql -p$ADMIN_PASSWORD -u root -D shrine_query_history < "$SHRINE_SRC_DIR/adapter/adapter-service/src/main/sql/adapter.sql"
mysql -p$ADMIN_PASSWORD -u root -D shrine_query_history < "$SHRINE_SRC_DIR/adapter/adapter-service/src/main/sql/mysql.ddl"
mysql -p$ADMIN_PASSWORD -u root -D shrine_query_history < "$SHRINE_SRC_DIR/hub/broadcaster-aggregator/src/main/sql/mysql.ddl"
mysql -p$ADMIN_PASSWORD -u root -D shrine_query_history < "$SHRINE_SRC_DIR/qep/service/src/main/sql/create_broadcaster_audit_table.sql"
mysql -p$ADMIN_PASSWORD -u root -D stewardDB < "$SHRINE_SRC_DIR/apps/steward-app/src/main/sql/mysql.ddl"
mysql -p$ADMIN_PASSWORD -u root -D adapterAuditDB < "$SHRINE_SRC_DIR/adapter/adapter-service/src/main/sql/mysql.ddl"
mysql -p$ADMIN_PASSWORD -u root -D qepAuditDB < "$SHRINE_SRC_DIR/qep/service/src/main/sql/mysql.ddl"
