#!/bin/bash
set -e

# expected envvar: I2B2_DATA_DIR, I2B2_DOMAIN_NAME

# ---------- CRC data ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Crcdata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2demodata
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2demodata
db.project=demo
EOL

ant -f data_build.xml create_crcdata_tables_release_1-7
ant -f data_build.xml create_procedures_release_1-7
#ant -f data_build.xml db_demodata_load_data


# ---------- Hive data ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Hivedata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2hive
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2hive
EOL

ant -f data_build.xml create_hivedata_tables_release_1-7
ant -f data_build.xml db_hivedata_load_data


# ---------- IM data ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Imdata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2imdata
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2imdata
db.project=demo
EOL

ant -f data_build.xml create_imdata_tables_release_1-7
#ant -f data_build.xml db_imdata_load_data


# ---------- Metadata ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Metadata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2metadata
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2metadata
db.project=demo
EOL

ant -f data_build.xml create_metadata_tables_release_1-7
#ant -f data_build.xml db_metadata_load_data


# ---------- PM data ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Pmdata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2pm
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2pm
EOL

ant -f data_build.xml create_pmdata_tables_release_1-7
ant -f data_build.xml create_triggers_release_1-7
ant -f data_build.xml db_pmdata_load_data


# ---------- Work data ----------
cd "$I2B2_DATA_DIR/edu.harvard.i2b2.data/Release_1-7/NewInstall/Workdata"

cat > db.properties <<EOL
db.type=postgresql
db.username=i2b2workdata
db.password=demouser
db.driver=org.postgresql.Driver
db.url=jdbc:postgresql://localhost/$I2B2_DOMAIN_NAME?searchpath=i2b2workdata
db.project=demo
EOL

ant -f data_build.xml create_workdata_tables_release_1-7
ant -f data_build.xml db_workdata_load_data
